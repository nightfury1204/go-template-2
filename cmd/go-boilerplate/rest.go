package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bitbucket.org/evaly/go-boilerplate/api"
	"bitbucket.org/evaly/go-boilerplate/config"
	"bitbucket.org/evaly/go-boilerplate/infra"
	infraMongo "bitbucket.org/evaly/go-boilerplate/infra/mongo"
	infraRedis "bitbucket.org/evaly/go-boilerplate/infra/redis"
	infraSentry "bitbucket.org/evaly/go-boilerplate/infra/sentry"
	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/repo"
	"bitbucket.org/evaly/go-boilerplate/service"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
)

// srvCmd is the serve sub command to start the api server
var srvCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves the api server",
	RunE:  serve,
}

func init() {
	srvCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func serve(cmd *cobra.Command, args []string) error {
	cfgApp := config.GetApp(cfgPath)
	cfgMongo := config.GetMongo(cfgPath)
	cfgRedis := config.GetRedis(cfgPath)
	cfgSentry := config.GetSentry(cfgPath)
	cfgDBTable := config.GetTable(cfgPath)

	ctx := context.Background()

	lgr := logger.DefaultOutStructLogger

	db, err := infraMongo.New(ctx, cfgMongo.URL, cfgMongo.DBName, cfgMongo.DBTimeOut)
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	kv, err := infraRedis.New(cfgRedis.URL, cfgRedis.RedisTimeOut, "go-boilerplate")
	if err != nil {
		return err
	}
	defer kv.Close()

	err = infraSentry.NewInit(cfgSentry.URL)
	if err != nil {
		return err
	}

	brandRepo := repo.NewBrand(cfgDBTable.BrandCollectionName, db)

	svc := service.NewBrand(brandRepo, kv, lgr)
	api.SetLogger(logger.DefaultOutLogger)

	errChan := make(chan error)
	go func() {
		if err := startHealthServer(cfgApp, db, kv); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := startApiServer(cfgApp, svc, lgr); err != nil {
			errChan <- err
		}
	}()
	return <-errChan

}

func startHealthServer(cfg *config.Application, db infra.DB, kv infra.KV) error {
	log.Println("startHealthServer")
	sc := api.NewSystemController(db, kv)
	api.NewSystemRouter(sc)
	r := chi.NewMux()
	r.Mount("/system/v1/", api.NewSystemRouter(sc))

	srvr := http.Server{
		Addr:    getAddressFromHostAndPort(cfg.Host, 3550),
		Handler: r,
		//ErrorLog: logger.DefaultErrLogger,
		//WriteTimeout: cfg.WriteTimeout,
		//ReadTimeout:  cfg.ReadTimeout,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	graceful := func() error {
		log.Println("To shutdown immedietly press again")

		return nil
	}

	errCh := make(chan error)
	forced := func() error {
		log.Println("Shutting down server forcefully")
		return nil
	}
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}

	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()

	return <-errCh
}

func startApiServer(cfg *config.Application, svc service.BrandService, lgr logger.StructLogger) error {
	brndsCtrl := api.NewBrandsController(svc)
	brndsCtrl.SetLogger(lgr)

	wrkrCtrl := api.NewWorkerController(svc)
	wrkrCtrl.SetLogger(lgr)

	r := chi.NewMux()
	r.Mount("/api/v1/public", api.NewRouter(brndsCtrl, wrkrCtrl))

	srvr := http.Server{
		Addr:    getAddressFromHostAndPort(cfg.Host, cfg.Port),
		Handler: r,
		//ErrorLog: logger.DefaultErrLogger,
		//WriteTimeout: cfg.WriteTimeout,
		//ReadTimeout:  cfg.ReadTimeout,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return ManageServer(&srvr, 30*time.Second)
}

func ManageServer(srvr *http.Server, gracePeriod time.Duration) error {
	errCh := make(chan error)

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt}

	graceful := func() error {
		log.Println("Suttingdown server gracefully with in", gracePeriod)
		log.Println("To shutdown immedietly press again")

		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		return srvr.Shutdown(ctx)
	}

	forced := func() error {
		log.Println("Shutting down server forcefully")
		return srvr.Close()
	}

	go func() {
		log.Println("Starting server on", srvr.Addr)
		if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()

	return <-errCh
}

// HandleSignals listen on the registered signals and fires the gracefulHandler for the
// first signal and the forceHandler (if any) for the next this function blocks and
// return any error that returned by any of the handlers first
func HandleSignals(sigs []os.Signal, gracefulHandler, forceHandler func() error) error {
	sigCh := make(chan os.Signal)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, sigs...)
	defer signal.Stop(sigCh)

	grace := true
	for {
		select {
		case err := <-errCh:
			return err
		case <-sigCh:
			if grace {
				grace = false
				go func() {
					errCh <- gracefulHandler()
				}()
			} else if forceHandler != nil {
				err := forceHandler()
				errCh <- err
			}
		}
	}
}

func getAddressFromHostAndPort(host string, port int) string {
	addr := host
	if port != 0 {
		addr = addr + ":" + strconv.Itoa(port)
	}
	return addr
}
