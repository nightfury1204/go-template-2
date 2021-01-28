package service

import (
	"context"

	"bitbucket.org/evaly/go-boilerplate/infra"
	"bitbucket.org/evaly/go-boilerplate/logger"
	"bitbucket.org/evaly/go-boilerplate/model"
	"bitbucket.org/evaly/go-boilerplate/repo"
	"bitbucket.org/evaly/go-boilerplate/utils"
)

// BrandService interface
type BrandService interface {
	ListBrand(ctx context.Context, pager utils.Pager) ([]model.BrandInfo, error)
	InsertNewBrandFromEvent(ctx context.Context, brand *model.BrandInfo) error
}

// Brand ...
type Brand struct {
	cache     infra.KV
	log       logger.StructLogger
	brandRepo repo.BrandRepo
}

// NewBrand ...
func NewBrand(brandRepo repo.BrandRepo, cache infra.KV, lgr logger.StructLogger) BrandService {
	return &Brand{
		cache:     cache,
		log:       lgr,
		brandRepo: brandRepo,
	}
}

// SetLogger ...
func (c *Brand) SetLogger(l logger.StructLogger) {
	c.log = l
}

// ListBrand ...
func (c *Brand) ListBrand(ctx context.Context, pager utils.Pager) ([]model.BrandInfo, error) {
	tid := utils.GetTracingID(ctx)
	brands, err := c.getListBrandsFromCache(pager)
	if err != nil {
		c.log.Errorf("ListBrands", tid, "failed to get from cache, %v", err)
	}
	if brands == nil {
		c.log.Println("ListBrands", tid, "listing product brands from database")
		brands, err = c.brandRepo.ListBrands(ctx, "", pager.Skip, pager.Limit)
		if err != nil {
			c.log.Errorln("ListBrands", tid, err.Error())
			return nil, err
		}
		if err := c.cacheListBrands(pager, brands); err != nil {
			c.log.Errorf("ListBrands", tid, "failed to cache, %v", err)
		}
	}

	c.log.Println("ListBrands", tid, "sent response successfully")
	return brands, nil
}
