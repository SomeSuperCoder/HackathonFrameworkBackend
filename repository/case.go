package repository

import (
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CaseRepo = GenericRepo[models.Case]

func NewCaseRepo(database *mongo.Database) *CaseRepo {
	return NewGenericRepo[models.Case](database, "cases")
}
