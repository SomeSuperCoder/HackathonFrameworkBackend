package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/internal"
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeamRepo struct {
	*GenericRepo[models.Team]
	database *mongo.Database
	Users    *mongo.Collection
}

func NewTeamRepo(database *mongo.Database) *TeamRepo {
	return &TeamRepo{
		database:    database,
		Users:       database.Collection("users"),
		GenericRepo: NewGenericRepo[models.Team](database, "teams"),
	}
}

func (r *TeamRepo) GetMembers(ctx context.Context, id bson.ObjectID) ([]models.User, error) {
	return FindWithFilter[models.User](ctx, r.Users, bson.M{
		"team": id,
	})
}

func (r *TeamRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	if err := Delete(ctx, r.Collection, id); err != nil {
		return err
	}

	_, err := r.Users.UpdateMany(ctx, bson.M{
		"team": id,
	}, bson.M{
		"$set": bson.M{
			"team": internal.UndefinedObjectID,
		},
	})
	return err
}
