package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
type PagedResponseBuilder[T any] = func(values []T, totalCount int64) any

type PagedFinder[T any] interface {
	FindPaged(ctx context.Context, page, limit int64) ([]T, int64, error)
}

func FindPaged[T any](w http.ResponseWriter, r *http.Request, repo PagedFinder[T], pagedResponseBuilder PagedResponseBuilder[T]) {
	// Get data
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Validate
	if page == "" {
		http.Error(w, "No page number provided", http.StatusBadRequest)
		return
	}
	if limit == "" {
		http.Error(w, "No limit number provided", http.StatusBadRequest)
		return
	}

	// Parse
	pageNumber, err := strconv.Atoi(page)
	if utils.CheckError(w, err, "Invalid page number", http.StatusBadRequest) {
		return
	}

	limitNumber, err := strconv.Atoi(limit)
	if utils.CheckError(w, err, "Invalid limit number", http.StatusBadRequest) {
		return
	}

	// Do work
	teams, totalCount, err := repo.FindPaged(r.Context(), int64(pageNumber), int64(limitNumber))
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	// Respond
	utils.RespondWithJSON(w, pagedResponseBuilder(teams, totalCount))
}

// ====================
type ValueGenerator[T any] = func() T

type Creatator[T any] interface {
	Create(ctx context.Context, value T) (bson.ObjectID, error)
}

func AdminOnlyCreate[T any, R any](w http.ResponseWriter, r *http.Request, repo Creatator[T], request R, valueGenerator ValueGenerator[T]) {
	Create(w, r, repo, request, func(w http.ResponseWriter, r *http.Request, id bson.ObjectID, userAuth *models.User) bool {
		return AdminCheck(w, r)
	}, valueGenerator)
}

func Create[T any, R any](w http.ResponseWriter, r *http.Request, repo Creatator[T], request R, accessChecker AccessChecker, valueGenerator ValueGenerator[T]) {
	if accessChecker(w, r, bson.NilObjectID, middleware.ExtractUserAuth(r)) {
		return
	}

	if DefaultParseAndValidate(w, r, request) {
		return
	}

	CreateInner(w, r, repo, valueGenerator())
}

func CreateInner[T any](w http.ResponseWriter, r *http.Request, repo Creatator[T], newValue T) {
	createdID, err := repo.Create(r.Context(), newValue)

	if utils.CheckError(w, err, "Failed to create", http.StatusInternalServerError) {
		return
	}

	fmt.Fprintln(w, createdID.Hex())
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
type Validator = func(w http.ResponseWriter, r *http.Request)

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
