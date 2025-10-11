package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GenericRepo[T any] struct {
	database   *mongo.Database
	Collection *mongo.Collection
}

func NewGenericRepo[T any](database *mongo.Database, collectionName string) *GenericRepo[T] {
	return &GenericRepo[T]{
		database:   database,
		Collection: database.Collection(collectionName),
	}
}

func (r *GenericRepo[T]) FindPaged(ctx context.Context, page, limit int64) ([]T, int64, error) {
	return FindPaged[T](ctx, r.Collection, page, limit)
}

func (r *GenericRepo[T]) Find(ctx context.Context) ([]T, error) {
	return Find[T](ctx, r.Collection)
}

func (r *GenericRepo[T]) Create(ctx context.Context, value *T) (bson.ObjectID, error) {
	return Create(ctx, r.Collection, value)
}

func (r *GenericRepo[T]) GetByID(ctx context.Context, id bson.ObjectID) (*T, error) {
	return GetByID[T](ctx, r.Collection, id)
}

func (r *GenericRepo[T]) Update(ctx context.Context, id bson.ObjectID, update any) error {
	return Update(ctx, r.Collection, id, update)
}

func (r *GenericRepo[T]) Delete(ctx context.Context, id bson.ObjectID) error {
	return Delete(ctx, r.Collection, id)
}

// =============================

func Find[T any](ctx context.Context, c *mongo.Collection) ([]T, error) {
	return FindWithFilter[T](ctx, c, struct{}{})
}

func FindWithFilter[T any](ctx context.Context, c *mongo.Collection, filter any) ([]T, error) {
	var values = []T{}

	// Init a cursor
	cursor, err := c.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Extract records
	err = cursor.All(ctx, &values)
	if err != nil {
		return nil, err
	}

	return values, err
}

func FindPaged[T any](ctx context.Context, c *mongo.Collection, page, limit int64) ([]T, int64, error) {
	var values = []T{}

	// Set pagination options
	skip := (page - 1) * limit
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(skip)
	opts.SetSort(bson.M{"created_at": -1})

	// Init a cursor
	cursor, err := c.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Extract records
	err = cursor.All(ctx, &values)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := c.CountDocuments(ctx, bson.M{})

	return values, count, err
}

func Create(ctx context.Context, c *mongo.Collection, value any) (bson.ObjectID, error) {
	res, err := c.InsertOne(ctx, value)
	objID, _ := res.InsertedID.(bson.ObjectID)
	return objID, err
}

func GetBy[T any](ctx context.Context, c *mongo.Collection, key string, value any) (*T, error) {
	var got T

	err := c.FindOne(ctx, bson.M{
		key: value,
	}, nil).Decode(&got)
	if err != nil {
		return nil, err
	}

	return &got, err
}

func GetByID[T any](ctx context.Context, c *mongo.Collection, id bson.ObjectID) (*T, error) {
	return GetBy[T](ctx, c, "_id", id)
}

func Update(ctx context.Context, c *mongo.Collection, id bson.ObjectID, update any) error {
	res := c.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
}

func Delete(ctx context.Context, c *mongo.Collection, id bson.ObjectID) error {
	_, err := c.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
