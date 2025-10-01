package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Case struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	ImageURI    string        `bson:"image_uri" json:"image_uri"`
}
