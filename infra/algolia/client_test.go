package algolia

import (
	"bitbucket.org/evaly/go-boilerplate/model"
	"fmt"
	"testing"
)

func TestAlgolia(t *testing.T) {
	client := New("FZFCF64BEM", "a81235e4ee92acd7459fb6b76577138e")
	client.InitializeIndex("products")
	err := client.CreateMany(getData())
	fmt.Println(err)
}

func TestAlgolia_DeleteMany(t *testing.T) {
	client := New("FZFCF64BEM", "a81235e4ee92acd7459fb6b76577138e")
	client.InitializeIndex("products")
	err := client.DeleteMany([]string{"1234"})
	fmt.Println(err)

}

func getData() []model.SearchEngineShopItem {
	data := make([]model.SearchEngineShopItem, 0)
	d := model.SearchEngineShopItem{
		ObjectID:        "1",
		ProductSlug:     "testproduct",
		ProductName:     "test product",
		ShopName:        "new test shop",
		Price:           10,
		DiscountedPrice: 9,
		ShopItemID:      1,
		Description:     "test description",
		ProductImage:    "testimageurl1",
		MinPrice:        10,
		MaxPrice:        8,
		BrandName:       "testbrand",
		CategoryName:    "test category",
		ColorVariants:   []string{"R", "G", "B"},
		Color:           "R",
	}
	data = append(data, d)
	return data
}
