package infra

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// DB interface wraps the databse
type DB interface {
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
	EnsureIndices(ctx context.Context, tab string, inds []DbIndex) error
	DropIndices(ctx context.Context, tab string, inds []DbIndex) error
	Insert(ctx context.Context, tab string, v interface{}) error
	InsertMany(ctx context.Context, tab string, v []interface{}) error
	List(ctx context.Context, tab string, filter DbQuery, skip, limit int64, v interface{}, sort ...interface{}) error
	FindOne(ctx context.Context, tab string, filter DbQuery, v interface{}, sort ...interface{}) error
	PartialUpdateMany(ctx context.Context, col string, filter DbQuery, data interface{}) error
	PartialUpdateManyByQuery(ctx context.Context, col string, filter DbQuery, query UnorderedDbQuery) error
	BulkUpdate(ctx context.Context, col string, models []mongo.WriteModel) error
	Aggregate(ctx context.Context, col string, q []DbQuery, v interface{}) error
	AggregateWithDiskUse(ctx context.Context, col string, q []DbQuery, v interface{}) error
	Distinct(ctx context.Context, col, field string, q DbQuery, v interface{}) error
	DeleteMany(ctx context.Context, col string, filter interface{}) error
}

// DbIndex holds database index
type DbIndex struct {
	Name   string
	Keys   []DbIndexKey
	Unique *bool
	Sparse *bool

	// If ExpireAfter is defined the server will periodically delete
	// documents with indexed time.Time older than the provided delta.
	ExpireAfter *time.Duration
}

type DbIndexKey struct {
	Key string
	Asc interface{}
}

// DbQuery holds a database query
type DbQuery bson.D
type UnorderedDbQuery bson.M

type BulkWriteModel mongo.WriteModel
