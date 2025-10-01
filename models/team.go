package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Grades map[bson.ObjectID]map[bson.ObjectID]uint16

func ValidateGrades(fl validator.FieldLevel) bool {
	// Get the field value and check if it's the correct type
	if nestedMap, ok := fl.Field().Interface().(Grades); ok {
		// Iterate through the outer map
		for _, innerMap := range nestedMap {
			// Iterate through each inner map
			for _, value := range innerMap {
				// Check if the integer value is within the required range
				if value > 10 {
					return false
				}
			}
		}
		return true // All values are valid
	}
	return false // Type assertion failed
}

type Team struct {
	ID              bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name            string        `bson:"name" json:"name"`
	Leader          bson.ObjectID `bson:"leader" json:"leader"`
	Repos           []string      `bson:"repos" json:"repos"`
	PresentationURI string        `bson:"presentationURI" json:"presentationURI"`
	Grades          Grades        `bson:"grades" json:"grades"`
	CratedAt        time.Time     `bson:"created_at" json:"created_at"`
}
