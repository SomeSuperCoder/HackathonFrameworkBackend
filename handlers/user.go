package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/internal/validators"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	Repo *repository.UserRepo
}

type UsersResponse struct {
	Users      []models.User `json:"users"`
	TotalCount int64         `json:"count"`
}

func (h *UserHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
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
	users, totalCount, err := h.Repo.FindPaged(r.Context(), int64(pageNumber), int64(limitNumber))
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	// Respond
	result := &UsersResponse{
		Users:      users,
		TotalCount: totalCount,
	}
	resultString, err := json.Marshal(result)
	if utils.CheckError(w, err, "Failed to serialize JSON", http.StatusInternalServerError) {
		return
	}

	fmt.Fprintln(w, string(resultString))
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	// Parse the user id
	parsedUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID provided", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByID(r.Context(), parsedUserID)
	getCommon(user, err, w)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	user, err := h.Repo.GetByUsername(r.Context(), username)
	getCommon(user, err, w)
}

func getCommon(user *models.User, err error, w http.ResponseWriter) {
	if utils.CheckGetFromDB(w, err) {
		return
	}

	serializedUser, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(serializedUser))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Load data
	id := r.PathValue("id")

	// Parse
	parsedId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID provided", http.StatusBadRequest)
		return
	}

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Parse
	var request struct {
		Name      string          `json:"name" bson:"name,omitempty" validate:"omitempty,self,min=1,max=20"`
		Birthdate time.Time       `json:"birthdate" bson:"birthdate,omitempty" validate:"omitempty,self"`
		Role      models.UserRole `json:"role" bson:"role,omitempty" validate:"omitempty,admin,oneof=0 1 2"`
		Team      bson.ObjectID   `json:"team" bson:"team,omitempty" validate:"omitempty,self"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	tv := validators.NewUserValidator(userAuth, parsedId)
	// Validate
	err = tv.ValidateRequest(&request)
	if utils.CheckJSONValidError(w, err) {
		return
	}

	// Do work
	err = h.Repo.Update(r.Context(), parsedId, request)
	if utils.CheckError(w, err, "Failed to update", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully updated")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Load data
	id := r.PathValue("id")

	// Parse
	parsedId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID provided", http.StatusBadRequest)
		return
	}

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Check access
	if userAuth.Role == models.Admin || parsedId == userAuth.ID {
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Do work
	err = h.Repo.Delete(r.Context(), parsedId)
	if utils.CheckError(w, err, "Failed to delete", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully deleted")
}
