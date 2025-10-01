package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Criterion struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Text string        `bson:"text" json:"text"`
}
