package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

type UserAuth struct {
	Username string
	UserID   bson.ObjectID
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	_, err := r.Users.InsertOne(ctx, user)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, userID bson.ObjectID) (*models.User, error) {
	return r.getCommon(ctx, bson.M{"_id": userID})
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.getCommon(ctx, bson.M{"username": username})
}

func (r *UserRepo) getCommon(ctx context.Context, filter bson.M) (*models.User, error) {
	opts := options.FindOne().SetProjection(bson.M{
		"sessions": 0,
	})

	var user models.User
	err := r.Users.FindOne(ctx, filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err

}
