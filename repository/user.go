package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	database *mongo.Database
	Users    *mongo.Collection
}

func NewUserRepo(database *mongo.Database) *UserRepo {
	return &UserRepo{
		database: database,
		Users:    database.Collection("users"),
	}
}

func (r *UserRepo) FindPaged(ctx context.Context, page, limit int64) ([]models.User, int64, error) {
	return FindPaged[models.User](ctx, r.Users, page, limit)
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (bson.ObjectID, error) {
	return Create(ctx, r.Users, user)
}

func (r *UserRepo) GetByID(ctx context.Context, userID bson.ObjectID) (*models.User, error) {
	return GetByID[models.User](ctx, r.Users, userID)
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return GetBy[models.User](ctx, r.Users, "username", username)
}

func (r *UserRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	res := r.Users.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
}

func (r *UserRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Users.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
