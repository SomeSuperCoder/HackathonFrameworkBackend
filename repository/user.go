package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepo struct {
	*GenericRepo[models.User]
}

func NewUserRepo(database *mongo.Database) *UserRepo {
	return &UserRepo{
		GenericRepo: NewGenericRepo[models.User](database, "users"),
	}
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return GetBy[models.User](ctx, r.Collection, "username", username)
}
