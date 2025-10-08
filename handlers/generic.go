package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/internal/validators"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GettableByID[T any] interface {
	GetByID(ctx context.Context, id bson.ObjectID) (T, error)
}

func GetByID[T any](w http.ResponseWriter, r *http.Request, repo GettableByID[T]) {
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	value, err := repo.GetByID(r.Context(), parsedId)
	if utils.CheckGetFromDB(w, err) {
		return
	}

	utils.RespondWithJSON(w, value)
}

// ====================

type Findable[T any] interface {
	Find(ctx context.Context) ([]T, error)
}

func Get[T any](w http.ResponseWriter, r *http.Request, repo Findable[T]) {
	cases, err := repo.Find(r.Context())
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	utils.RespondWithJSON(w, cases)
}

// ===================================================
// Helpers
// ===================================================
func ParseAndValidate(w http.ResponseWriter, r *http.Request, validator validators.Validator, request any) bool {
	err := json.NewDecoder(r.Body).Decode(request)
	if utils.CheckJSONError(w, err) {
		return true
	}

	err = validator.ValidateRequest(request)
	if utils.CheckJSONValidError(w, err) {
		return true
	}

	return false
}

func DefaultParseAndValidate(w http.ResponseWriter, r *http.Request, request any) bool {
	return ParseAndValidate(w, r, validators.NewAccessValidator(middleware.ExtractUserAuth(r)), request)
}

func AdminCheck(w http.ResponseWriter, r *http.Request) bool {
	if middleware.ExtractUserAuth(r).Role != models.Admin {
		http.Error(w, "Access denied: only the admin can perform this operation", http.StatusForbidden)
		return true
	}

	return false
}
