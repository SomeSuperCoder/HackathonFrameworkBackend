package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Grades map[bson.ObjectID]map[bson.ObjectID]uint16

type Team struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name            string        `bson:"name" json:"name"`
	Leader          bson.ObjectID `bson:"leader" json:"leader"`
	Repos           []string      `bson:"repos" json:"repos"`
	PresentationURI string        `bson:"presentationURI" json:"presentationURI"`
	Grades          Grades        `bson:"grades" json:"grades"`
	CratedAt        time.Time     `bson:"created_at" json:"created_at"`
}
