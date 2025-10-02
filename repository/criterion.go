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
	var criteria = []models.Criterion{}

	// Init a cursor
	cursor, err := r.Criteria.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Extract records
	err = cursor.All(ctx, &criteria)
	if err != nil {
		return nil, err
	}

	return criteria, err
}

func (r *CriterionRepo) Create(ctx context.Context, criterion *models.Criterion) error {
	_, err := r.Criteria.InsertOne(ctx, criterion)
	return err
}

func (r *CriterionRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Criterion, error) {
	var criterion models.Criterion

	err := r.Criteria.FindOne(ctx, bson.M{
		"_id": id,
	}, nil).Decode(&criterion)
	if err != nil {
		return nil, err
	}

	return &criterion, err
}

func (r *CriterionRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	res := r.Criteria.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
}

func (r *CriterionRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Criteria.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
