package service

import (
	"context"

	"bitbucket.org/evaly/go-boilerplate/model"
	"bitbucket.org/evaly/go-boilerplate/utils"
)

func (c *Brand) InsertNewBrandFromEvent(ctx context.Context, brand *model.BrandInfo) error {
	tid := utils.GetTracingID(ctx)
	c.log.Println("InsertNewBrandFromEvent", tid, "inserting shop into database")
	if err := c.brandRepo.Create(ctx, brand); err != nil {
		c.log.Errorln("InsertNewBrandFromEvent", tid, err.Error())
		return err
	}
	c.log.Println("InsertNewBrandFromEvent", tid, "sending response")

	return nil
}
