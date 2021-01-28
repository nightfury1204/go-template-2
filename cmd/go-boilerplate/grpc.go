package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/evaly/go-boilerplate/config"
	infraMongo "bitbucket.org/evaly/go-boilerplate/infra/mongo"
	infraRedis "bitbucket.org/evaly/go-boilerplate/infra/redis"
	infraSentry "bitbucket.org/evaly/go-boilerplate/infra/sentry"
	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/repo"
	"bitbucket.org/evaly/go-boilerplate/rpcs/pb"
	"bitbucket.org/evaly/go-boilerplate/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpcs "bitbucket.org/evaly/go-boilerplate/rpcs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// srvCmd is the serve sub command to start the api server
var grpcSrvCmd = &cobra.Command{
	Use:   "serve-grpc",
	Short: "serve serves the grpc server",
	RunE:  serveGrpc,
}

func init() {
	grpcSrvCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func serveGrpc(cmd *cobra.Command, args []string) error {
	cfgApp := config.GetApp(cfgPath)
	cfgMongo := config.GetMongo(cfgPath)
	cfgRedis := config.GetRedis(cfgPath)
	cfgSentry := config.GetSentry(cfgPath)
	cfgDBTable := config.GetTable(cfgPath)

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

	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Println("started server: ", fmt.Sprintf("%s:%d", cfgApp.Host, cfgApp.Port))
		if err := grpcSrvr.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Printf("Fatal error: %v\n", err)
	case <-stopChan:
		log.Println("initiating graceful shut down")
		grpcSrvr.GracefulStop()
		log.Println("graceful shut down done")
	}

	return nil
}
