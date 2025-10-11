package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CriterionRepo struct {
	database *mongo.Database
	Criteria *mongo.Collection
}

func NewCriterionRepo(database *mongo.Database) *CriterionRepo {
	return &CriterionRepo{
		database: database,
		Criteria: database.Collection("criteria"),
	}
}

func (r *CriterionRepo) Find(ctx context.Context) ([]models.Criterion, error) {
	return Find[models.Criterion](ctx, r.Criteria)
}

func (r *CriterionRepo) Create(ctx context.Context, criterion *models.Criterion) (bson.ObjectID, error) {
	return Create(ctx, r.Criteria, criterion)
}

func (r *CriterionRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Criterion, error) {
	return GetByID[models.Criterion](ctx, r.Criteria, id)
}

func (r *CriterionRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	return Update(ctx, r.Criteria, id, update)
}

func (r *CriterionRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	return Delete(ctx, r.Criteria, id)
}
