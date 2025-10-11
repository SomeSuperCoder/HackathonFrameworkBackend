package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
