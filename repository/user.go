package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepo struct {
	Database *mongo.Database
}

type UserAuth struct {
	Username string
	UserID   bson.ObjectID
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.Database.Collection("users").InsertOne(ctx, user)
	return err
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID bson.ObjectID) (*models.User, error) {
	return r.getUserCommon(ctx, bson.M{"_id": userID})
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.getUserCommon(ctx, bson.M{"username": username})
}

func (r *UserRepo) getUserCommon(ctx context.Context, filter bson.M) (*models.User, error) {
	opts := options.FindOne().SetProjection(bson.M{
		"sessions": 0,
	})

	var user models.User
	err := r.Database.Collection("users").FindOne(ctx, filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err

}
