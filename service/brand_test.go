package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"

	"bitbucket.org/evaly/go-boilerplate/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func dbConn() *mongo.Client {
	//uri := "mongodb+srv://catalogUser:9guEEVpPpyI2uzrK@catalog.bs48e.mongodb.net/catalog?retryWrites=true&w=majority"
	uri := "mongodb://adminUser:1qazZAQ!@13.251.114.75:27017/catalog?authSource=admin"
	//uri := "mongodb://root:secret@localhost:27017/catalog?authSource=admin"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Print(err)
	}
	return client

}

func TestInsertBrand(t *testing.T) {
	c := dbConn()
	brandCollection := c.Database("catalog").Collection("brand")
	_, err := brandCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	brandArray := getBrands()
	for _, prd := range brandArray {
		res, err := brandCollection.InsertOne(context.Background(), prd)
		if err != nil {
			log.Println(err)
			continue
		}
		t.Log("inserted in brand with id: ", res.InsertedID)
	}
}

func getBrands() []model.BrandInfo {
	brands := make([]model.BrandInfo, 0)
	for i := 0; i < 100; i++ {
		bid := strconv.Itoa(i)
		bi := model.BrandInfo{
			ID:          int64(i + 1),
			Name:        "Brand-" + bid,
			Approved:    true,
			Slug:        "brand_" + bid,
			Description: "Brand Description",
			BrandType:   "brandtype",
			ImageURL:    "image" + bid,
			Status:      model.StatusActive,
			BrandScore:  rand.Float64(),
		}
		brands = append(brands, bi)
	}
	return brands
}
