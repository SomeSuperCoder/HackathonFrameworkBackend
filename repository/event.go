package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventRepo struct {
	database *mongo.Database
	Events   *mongo.Collection
}

func NewEventRepo(database *mongo.Database) *EventRepo {
	return &EventRepo{
		database: database,
		Events:   database.Collection("events"),
	}
}

func (r *EventRepo) Find(ctx context.Context) ([]models.Event, error) {
	return Find[models.Event](ctx, r.Events)
}

func (r *EventRepo) Create(ctx context.Context, event *models.Event) error {
	_, err := r.Events.InsertOne(ctx, event)
	return err
}

func (r *EventRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Event, error) {
	var event models.Event

	err := r.Events.FindOne(ctx, bson.M{
		"_id": id,
	}, nil).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, err
}

func (r *EventRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	res := r.Events.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
}

func (r *EventRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Events.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
