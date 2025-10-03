package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRole int

const (
	Participant UserRole = iota
	Judge
	Admin
)

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	// Personal data
	Name      string        `bson:"name" json:"name"`
	Birthdate time.Time     `bson:"birthdate" json:"birthdate"`
	Role      UserRole      `bson:"role" json:"role"`
	Team      bson.ObjectID `bson:"team" json:"team"`
	// TG related
	Username string `bson:"username" json:"username"`
	ChatID   int64  `bson:"chat_id" json:"chat_id"`
}
