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

func (r *EventRepo) Create(ctx context.Context, event *models.Event) (bson.ObjectID, error) {
	return Create(ctx, r.Events, event)
}

func (r *EventRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Event, error) {
	return GetByID[models.Event](ctx, r.Events, id)
}

func (r *EventRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	return Update(ctx, r.Events, id, update)
}

func (r *EventRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	return Delete(ctx, r.Events, id)
}
