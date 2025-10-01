package repository

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	var messages = []models.Team{}

	// Set pagination options
	skip := (page - 1) * limit
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(skip)
	opts.SetSort(bson.M{"created_at": -1})

	// Init a cursor
	cursor, err := r.Teams.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Extract messages
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, 0, err
	}

	// Get total message count
	count, err := r.Teams.CountDocuments(ctx, bson.M{})

	return messages, count, err
}

func (r *TeamRepo) Create(ctx context.Context, team *models.Team) error {
	_, err := r.Teams.InsertOne(ctx, team)
	return err
}

func (r *TeamRepo) GetByID(ctx context.Context, id bson.ObjectID) (*models.Team, error) {
	var team models.Team
	err := r.Teams.FindOne(ctx, bson.M{
		"_id": id,
	}, nil).Decode(&team)
	if err != nil {
		return nil, err
	}

	return &team, err
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
	}, update)

	return res.Err()
}

func (r *TeamRepo) DeleteTeam(ctx context.Context, id bson.ObjectID) error {
	_, err := r.Teams.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	_, err = r.Users.UpdateMany(ctx, bson.M{
		"team": id,
	}, bson.M{
		"team": bson.NilObjectID,
	})
	if err != nil {
		return err
	}

	return nil
}
