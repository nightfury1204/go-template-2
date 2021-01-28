package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"bitbucket.org/evaly/go-boilerplate/config"
	infraMongo "bitbucket.org/evaly/go-boilerplate/infra/mongo"
	infraRedis "bitbucket.org/evaly/go-boilerplate/infra/redis"
	infraSentry "bitbucket.org/evaly/go-boilerplate/infra/sentry"
	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/repo"
	"bitbucket.org/evaly/go-boilerplate/rpcrestproxy"
	"bitbucket.org/evaly/go-boilerplate/rpcrestproxy/handler"
	rpcs "bitbucket.org/evaly/go-boilerplate/rpcs"
	"bitbucket.org/evaly/go-boilerplate/rpcs/pb"
	"bitbucket.org/evaly/go-boilerplate/service"
	"github.com/go-chi/chi"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// srvCmd is the serve sub command to start the api server
var grpcRestSrvCmd = &cobra.Command{
	Use:   "serve-grpc-rest",
	Short: "serve serves the grpc rest server",
	RunE:  serveGrpcRest,
}

func init() {
	grpcRestSrvCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func serveGrpcForRest(errChan chan error) error {
	cfgApp := config.GetApp(cfgPath)
	cfgMongo := config.GetMongo(cfgPath)
	cfgRedis := config.GetRedis(cfgPath)
	cfgSentry := config.GetSentry(cfgPath)
	cfgDBTable := config.GetTable(cfgPath)

	fmt.Println(cfgApp, cfgMongo, cfgRedis, cfgSentry, cfgDBTable)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfgApp.Host, cfgApp.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	customFunc := func(p interface{}) (err error) {
		log.Println("panic triggered: ", p)
		return status.Errorf(codes.Unknown, "something went wrong")
	}

	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	grpcSrvr := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			//AuthUnaryInterceptor,
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			//AuthstreamInterceptor,
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)

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

	brndHndlr := rpcs.NewBrandServer(svc)

	// registering grpc handler
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcSrvr, healthcheck)

	pb.RegisterBrandServiceServer(grpcSrvr, brndHndlr)

	errChanLc := make(chan error)
	// stopChan := make(chan os.Signal)
	// signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Println("started server: ", fmt.Sprintf("%s:%d", cfgApp.Host, cfgApp.Port))
		if err := grpcSrvr.Serve(lis); err != nil {
			errChanLc <- err
		}
	}()

	err = <-errChanLc
	errChan <- err
	log.Printf("Fatal error: %v\n", err)

	return nil
}

func serveGrpcRestAPIServer(cfg *config.Application, clnt *rpcrestproxy.GoBoilerplateClients, lgr logger.StructLogger) error {
	brndsCtrl := handler.NewBrandHandler(clnt, lgr)

	r := chi.NewMux()
	r.Mount("/api/v1/public", handler.GetRouter(brndsCtrl))

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

func serveGrpcRestServer(errChan chan error) error {
	cfgApp := config.GetApp(cfgPath)
	cfgMongo := config.GetMongo(cfgPath)
	cfgRedis := config.GetRedis(cfgPath)
	cfgSentry := config.GetSentry(cfgPath)
	cfgRPC := config.GetRPC(cfgPath)

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

	clnt, err := rpcrestproxy.NewGoBoilerplateClients(cfgRPC.URL, lgr)
	if err != nil {
		return err
	}

	errChanLc := make(chan error)

	go func() {
		if err := startHealthServer(cfgApp, db, kv); err != nil {
			errChanLc <- err
		}
	}()

	go func() {
		if err := serveGrpcRestAPIServer(cfgApp, clnt, lgr); err != nil {
			errChanLc <- err
		}
	}()

	err = <-errChanLc
	errChanLc <- err

	return err
}

func serveGrpcRest(cmd *cobra.Command, args []string) error {
	errChan := make(chan error)

	go serveGrpcForRest(errChan)

	go serveGrpcRestServer(errChan)

	return <-errChan
}
