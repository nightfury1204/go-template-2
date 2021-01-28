package rpcrestproxy

import (
	"log"
	"net/http"

	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/rpcs/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

// GoBoilerplateClients holds go boilerplate app grpc client
type GoBoilerplateClients struct {
	BrandServiceClient pb.BrandServiceClient
	lgr                logger.StructLogger
}

// NewGoBoilerplateClients returns new grpc  client
func NewGoBoilerplateClients(rpcServerURL string, lgr logger.StructLogger) (*GoBoilerplateClients, error) {
	log.Println("connecting grpc server, grpc dialing to: ", rpcServerURL)
	con, err := grpc.Dial(
		rpcServerURL, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				grpc_retry.UnaryClientInterceptor(),
			),
		),
		grpc.WithChainStreamInterceptor(
			grpc_middleware.ChainStreamClient(
				grpc_retry.StreamClientInterceptor(),
			),
		),
	)
	if err != nil {
		log.Fatalln("failed to connect order service")
		return nil, err
	}

	return &GoBoilerplateClients{
		BrandServiceClient: pb.NewBrandServiceClient(con),
		lgr:                lgr,
	}, nil
}

// HTTPStatusFromCode returns http status code from grpc code
func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}
