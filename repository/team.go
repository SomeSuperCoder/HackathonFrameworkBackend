package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/internal"
	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeamRepo struct {
	database *mongo.Database
	Teams    *mongo.Collection
	Users    *mongo.Collection
}

func NewTeamRepo(database *mongo.Database) *TeamRepo {
	return &TeamRepo{
		database: database,
		Teams:    database.Collection("teams"),
		Users:    database.Collection("users"),
	}
}

func (r *TeamRepo) FindPaged(ctx context.Context, page, limit int64) ([]models.Team, int64, error) {
	return FindPaged[models.Team](ctx, r.Teams, page, limit)
}

func (r *TeamRepo) Create(ctx context.Context, team *models.Team) (bson.ObjectID, error) {
	return Create(ctx, r.Teams, team)
}

func (r *TeamRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Team, error) {
	return GetByID[models.Team](ctx, r.Teams, id)
}

func (r *TeamRepo) GetMembers(ctx context.Context, id bson.ObjectID) ([]models.User, error) {
	var members = []models.User{}
	cursor, err := r.Users.Find(ctx, bson.M{
		"team": id,
	}, nil)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *TeamRepo) Update(ctx context.Context, id bson.ObjectID, update any) error {
	res := r.Teams.FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": update,
	})

	return res.Err()
}

func (r *TeamRepo) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Teams.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	_, err = r.Users.UpdateMany(ctx, bson.M{
		"team": id,
	}, bson.M{
		"$set": bson.M{
			"team": internal.UndefinedObjectID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
