package repo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// ToBsonMDoc marshals to bson.M
func ToBsonMDoc(v interface{}) (*bson.M, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	var doc bson.M
	err = bson.Unmarshal(data, &doc)
	return &doc, err
}

// ToBsonDDoc marshals to bson.D
func ToBsonDDoc(v interface{}) (bson.D, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	var doc bson.D
	err = bson.Unmarshal(data, &doc)
	return doc, err
}

// Repo defines repository
type Repo interface {
	EnsureIndices() error
	DropIndices() error
}
