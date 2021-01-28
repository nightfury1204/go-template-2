package response

import (
	"bitbucket.org/evaly/go-boilerplate/model"
	"bitbucket.org/evaly/go-boilerplate/rpcs/pb"
)

// ToPbBrand ...
func ToPbBrand(req *model.BrandInfo) *pb.Brand {
	return &pb.Brand{}
}

// ToPbBrands ...
func ToPbBrands(req []model.BrandInfo) *pb.ListBrandRes {
	pi := []*pb.Brand{}

	for _, v := range req {
		pi = append(pi, ToPbBrand(&v))
	}

	return &pb.ListBrandRes{
		Brands: pi,
	}
}
