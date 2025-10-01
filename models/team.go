package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Team struct {
	ID              bson.ObjectID           `bson:"_id,omitempty" json:"_id"`
	Leader          bson.ObjectID           `bson:"leader" json:"leader"`
	Repos           []string                `bson:"repos" json:"repos"`
	PresentationURI string                  `bson:"presentationURI" json:"presentationURI"`
	Grades          map[bson.ObjectID][]int `bson:"grades" json:"grades"`
	CratedAt        time.Time               `bson:"created_at" json:"created_at"`
}
