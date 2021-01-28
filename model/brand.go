package model

// Status defines status
type Status string

const (
	StatusActive   Status = "active"
	StatusArchived Status = "archived"
)

// BrandInfo defines brand model
type BrandInfo struct {
	ID          int64    `json:"id,omitempty" bson:"id"`
	Name        string   `json:"name,omitempty" bson:"name"`
	Approved    bool     `json:"approved,omitempty" bson:"approved"`
	Slug        string   `json:"slug,omitempty" bson:"slug"`
	Description string   `json:"description,omitempty" bson:"description"`
	BrandType   string   `json:"brand_type,omitempty" bson:"brand_type"`
	Categories  []string `json:"categories,omitempty" bson:"categories"`
	ImageURL    string   `json:"image_url,omitempty" bson:"image_url"`
	Status      Status   `json:"status,omitempty" bson:"status"`
	BrandScore  float64  `json:"brand_score,omitempty" bson:"brand_score"`
	Version     int64    `json:"version" bson:"version"`
}
