package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Find[T any](ctx context.Context, c *mongo.Collection) ([]T, error) {
	var values = []T{}

	// Init a cursor
	cursor, err := c.Find(ctx, bson.M{}, nil)
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
	var teams = []T{}

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
	err = cursor.All(ctx, &teams)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := c.CountDocuments(ctx, bson.M{})

	return teams, count, err
}
