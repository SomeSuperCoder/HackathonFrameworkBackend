package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Event struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Time        time.Time     `bson:"time" json:"time"`
}
