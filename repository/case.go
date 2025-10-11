package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CaseRepo struct {
	database *mongo.Database
	Cases    *mongo.Collection
}

func NewCaseRepo(database *mongo.Database) *CaseRepo {
	return &CaseRepo{
		database: database,
		Cases:    database.Collection("cases"),
	}
}

func (r *CaseRepo) Find(ctx context.Context) ([]models.Case, error) {
	return Find[models.Case](ctx, r.Cases)
}

func (r *CaseRepo) Create(ctx context.Context, case_ *models.Case) (bson.ObjectID, error) {
	return Create(ctx, r.Cases, case_)
}

func (r *CaseRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Case, error) {
	return GetByID[models.Case](ctx, r.Cases, id)
}

func (r *CaseRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	return Update(ctx, r.Cases, id, update)
}

func (r *CaseRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Cases.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
