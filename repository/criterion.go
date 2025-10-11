package repository

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CriterionRepo = GenericRepo[models.Criterion]

func NewCriterionRepo(database *mongo.Database) *CriterionRepo {
	return NewGenericRepo[models.Criterion](database, "criteria")
}
