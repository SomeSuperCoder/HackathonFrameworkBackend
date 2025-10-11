package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/internal/validators"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GetteryID[T any] interface {
	GetByID(ctx context.Context, id bson.ObjectID) (T, error)
}

func GetByID[T any](w http.ResponseWriter, r *http.Request, repo GetteryID[T]) {
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
type Finder[T any] interface {
	Find(ctx context.Context) ([]T, error)
}

func Get[T any](w http.ResponseWriter, r *http.Request, repo Finder[T]) {
	cases, err := repo.Find(r.Context())
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	utils.RespondWithJSON(w, cases)
}

// ====================
type Creatator[T any] interface {
	Create(ctx context.Context, value T) error
}

func Create[T any](w http.ResponseWriter, r *http.Request, repo Creatator[T], request T) {
	err := repo.Create(r.Context(), request)

	if utils.CheckError(w, err, "Failed to create", http.StatusInternalServerError) {
		return
	}

	fmt.Fprintf(w, "Successfully created")
}

// ====================
type Updater interface {
	Update(ctx context.Context, id bson.ObjectID, update any) error
}

func Update[T any](w http.ResponseWriter, r *http.Request, repo Updater, request T) {
	// Load data
	var id bson.ObjectID
	var exit bool
	if id, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Parse
	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	UpdateInner(w, r, repo, id, request)
}

// ====================
type Deleter interface {
	Delete(ctx context.Context, id bson.ObjectID) error
}

func AdminOnlyDelete(w http.ResponseWriter, r *http.Request, repo Deleter) {
	Delete(w, r, repo, func(w http.ResponseWriter, r *http.Request, id bson.ObjectID, userAuth *models.User) bool {
		return AdminCheck(w, r)
	})
}

func Delete(w http.ResponseWriter, r *http.Request, repo Deleter, accessChecker AccessChecker) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	if accessChecker(w, r, parsedId, userAuth) {
		return
	}

	// Do work
	err := repo.Delete(r.Context(), parsedId)
	if utils.CheckError(w, err, "Failed to delete", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully deleted")
}

// ===================================================
// Helpers
// ===================================================
type AccessChecker = func(w http.ResponseWriter, r *http.Request, id bson.ObjectID, userAuth *models.User) bool

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

func UpdateInner(w http.ResponseWriter, r *http.Request, repo Updater, id bson.ObjectID, request any) {
	// Do work
	err := repo.Update(r.Context(), id, request)
	if utils.CheckError(w, err, "Failed to update", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully updated")
}

func AdminCheck(w http.ResponseWriter, r *http.Request) bool {
	if middleware.ExtractUserAuth(r).Role != models.Admin {
		http.Error(w, "Access denied: only the admin can perform this operation", http.StatusForbidden)
		return true
	}

	return false
}
