package handlers

import (
	"context"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GettableByID[T any] interface {
	GetByID(ctx context.Context, id bson.ObjectID) (T, error)
}

func GetByID[T any](w http.ResponseWriter, r *http.Request, repo GettableByID[T]) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Do work
	value, err := repo.GetByID(r.Context(), parsedId)
	if utils.CheckGetFromDB(w, err) {
		return
	}

	// Respond
	utils.RespondWithJSON(w, value)
}

// ====================

type Findable[T any] interface {
	Find(ctx context.Context) ([]T, error)
}

func Get[T any](w http.ResponseWriter, r *http.Request, repo Findable[T]) {
	//Do work
	cases, err := repo.Find(r.Context())
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	// Respond
	utils.RespondWithJSON(w, cases)
}
