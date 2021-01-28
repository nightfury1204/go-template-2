package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/evaly/go-boilerplate/infra"
	"bitbucket.org/evaly/go-boilerplate/model"
	"bitbucket.org/evaly/go-boilerplate/utils"
)

const (
	tableBrandCache    = "brands"
	keyListBrandsCache = "list_brands"

	defaultCacheDuration        = 10 * time.Minute
	defaultGeneralCacheDuration = 15 * time.Minute
	defaultLongCacheDuration    = 1 * time.Hour
)

func buildKeyFromPager(cacheKey string, pager utils.Pager) string {
	var key strings.Builder
	key.WriteString(fmt.Sprintf("%s_%d_%d", cacheKey, pager.Skip, pager.Limit))
	return key.String()
}

func (c *Brand) getListBrandsFromCache(pager utils.Pager) ([]model.BrandInfo, error) {
	key := buildKeyFromPager(keyListBrandsCache, pager)
	resp := make([]model.BrandInfo, 0)

	if err := c.cache.Get(context.Background(), tableBrandCache, key, &resp); err != nil {
		if err == infra.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return resp, nil
}

func (c *Brand) cacheListBrands(pager utils.Pager, data []model.BrandInfo) error {
	key := buildKeyFromPager(keyListBrandsCache, pager)
	return c.cache.PutEx(context.Background(), tableBrandCache, key, data, defaultCacheDuration)
}
