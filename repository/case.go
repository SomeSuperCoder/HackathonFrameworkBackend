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
	var cases = []models.Case{}

	// Init a cursor
	cursor, err := r.Cases.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Extract records
	err = cursor.All(ctx, &cases)
	if err != nil {
		return nil, err
	}

	return cases, err
}

func (r *CaseRepo) Create(ctx context.Context, case_ *models.Case) error {
	_, err := r.Cases.InsertOne(ctx, case_)
	return err
}

func (r *CaseRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Case, error) {
	var case_ models.Case

	err := r.Cases.FindOne(ctx, bson.M{
		"_id": id,
	}, nil).Decode(&case_)
	if err != nil {
		return nil, err
	}

	return &case_, err
}

func (r *CaseRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	res := r.Cases.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
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
