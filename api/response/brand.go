package response

import "bitbucket.org/evaly/go-boilerplate/model"

type BrandInfoResp struct {
	ID          int64        `json:"id,omitempty" bson:"id"`
	Name        string       `json:"name,omitempty" bson:"name"`
	Approved    bool         `json:"approved,omitempty" bson:"approved"`
	Slug        string       `json:"slug,omitempty" bson:"slug"`
	Description string       `json:"description,omitempty" bson:"description"`
	BrandType   string       `json:"brand_type,omitempty" bson:"brand_type"`
	ImageURL    string       `json:"image_url,omitempty" bson:"image_url"`
	Status      model.Status `json:"status,omitempty" bson:"status"`
	BrandScore  float64      `json:"brand_score,omitempty" bson:"brand_score"`
	Version     int64        `json:"version" bson:"version"`
}

func ToBrandInfoResp(brands []model.BrandInfo) []BrandInfoResp {
	resp := make([]BrandInfoResp, 0)
	for _, brand := range brands {
		resp = append(resp, BrandInfoResp{
			ID:          brand.ID,
			Name:        brand.Name,
			Approved:    brand.Approved,
			Slug:        brand.Slug,
			Description: brand.Description,
			BrandType:   brand.BrandType,
			ImageURL:    brand.ImageURL,
			Status:      brand.Status,
			BrandScore:  brand.BrandScore,
			Version:     brand.Version,
		})
	}
	return resp
}
