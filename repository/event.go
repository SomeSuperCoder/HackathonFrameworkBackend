package repository

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventRepo = GenericRepo[models.Event]

func NewEventRepo(database *mongo.Database) *EventRepo {
	return NewGenericRepo[models.Event](database, "events")
}
