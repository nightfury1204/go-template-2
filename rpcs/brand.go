package rpc

import (
	"context"

	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/rpcs/pb"
	"bitbucket.org/evaly/go-boilerplate/rpcs/response"
	"bitbucket.org/evaly/go-boilerplate/service"
	"bitbucket.org/evaly/go-boilerplate/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BrandsServer ...
type BrandsServer struct {
	svc service.BrandService
	lgr logger.StructLogger
}

// NewBrandServer ...
func NewBrandServer(svc service.BrandService) *BrandsServer {
	return &BrandsServer{
		svc: svc,
	}
}

func validateListBrandReq(req *pb.ListBrandReq) (bool, error) {
	ve := &errdetails.BadRequest{}

	if len(ve.FieldViolations) > 0 {
		st := status.New(codes.InvalidArgument, "invalid request body")
		st, err := st.WithDetails(ve)
		if err != nil {
			return false, status.Errorf(codes.InvalidArgument, "invalid fields")
		}
		return false, st.Err()
	}

	return true, nil
}

// ListBrand ...
func (brndSrvr *BrandsServer) ListBrand(ctx context.Context, req *pb.ListBrandReq) (*pb.ListBrandRes, error) {
	ok, err := validateListBrandReq(req)
	if !ok {
		return nil, err
	}

	res, err := brndSrvr.svc.ListBrand(context.Background(), utils.Pager{Skip: 0, Limit: 24})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}

	return response.ToPbBrands(res), nil
}
