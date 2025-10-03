package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/internal/validators"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeamHandler struct {
	Repo *repository.TeamRepo
}

type TeamsResponse struct {
	Teams      []models.Team `json:"teams"`
	TotalCount int64         `json:"count"`
}

func (h *TeamHandler) GetPaged(w http.ResponseWriter, r *http.Request) {
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
	teams, totalCount, err := h.Repo.FindPaged(r.Context(), int64(pageNumber), int64(limitNumber))
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	// Respond
	result := &TeamsResponse{
		Teams:      teams,
		TotalCount: totalCount,
	}
	resultString, err := json.Marshal(result)
	if utils.CheckError(w, err, "Failed to serialize JSON", http.StatusInternalServerError) {
		return
	}

	fmt.Fprintln(w, string(resultString))
}

func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Do work
	team, err := h.Repo.GetByID(r.Context(), parsedId)
	if utils.CheckGetFromDB(w, err) {
		return
	}

	// Respond
	serialized, err := json.Marshal(&team)
	if utils.CheckJSONError(w, err) {
		return
	}

	fmt.Fprintln(w, string(serialized))
}

func (h *TeamHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Do work
	members, err := h.Repo.GetMembers(r.Context(), parsedId)
	if utils.CheckError(w, err, "Failed to get members", http.StatusInternalServerError) {
		return
	}

	// Respond
	serialized, err := json.Marshal(&members)
	if utils.CheckJSONError(w, err) {
		return
	}

	fmt.Fprintln(w, string(serialized))
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Check access
	if userAuth.Team != bson.NilObjectID {
		http.Error(w, "Access denied: you already are part of a team", http.StatusForbidden)
		return
	}

	// Parse
	var request struct {
		Name string `json:"name" bson:"name,omitempty" validate:"omitempty,min=1,max=40"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	// Validate
	tv := validators.NewTeamValidator(userAuth, nil)
	err = tv.ValidateRequest(&request)
	if utils.CheckJSONValidError(w, err) {
		return
	}

	// Do work
	err = h.Repo.Create(r.Context(), &models.Team{
		Name:            request.Name,
		Leader:          userAuth.ID,
		Repos:           make([]string, 0),
		PresentationURI: "",
		Grades:          make(models.Grades),
	})

	if utils.CheckError(w, err, "Failed to create", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully created")
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Check if team exists
	team, err := h.Repo.GetByID(r.Context(), parsedId)
	if errors.Is(err, mongo.ErrNoDocuments) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	} else if utils.CheckError(w, err, "Failed to get team from DB", http.StatusInternalServerError) {
		return
	}

	// Parse
	var request struct {
		Name            string        `json:"name" bson:"name,omitempty" validate:"omitempty,owner,min=1,max=40"`
		Leader          bson.ObjectID `json:"leader" bson:"leader,omitempty" validate:"omitempty,owner"`
		Repos           []string      `json:"repos" bson:"repos,omitempty" validate:"omitempty,owner,dive,url"`
		PresentationURI string        `json:"presentation_uri" bson:"presentation_uri,omitempty" validate:"omitempty,owner,url"`
		Grades          models.Grades `json:"grades" bson:"grades,omitempty" validate:"omitempty,judge,grades"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	tv := validators.NewTeamValidator(userAuth, team)
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

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Check access
	if userAuth.Role == models.Admin || parsedId == userAuth.Team {
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Do work
	err := h.Repo.Delete(r.Context(), parsedId)
	if utils.CheckError(w, err, "Failed to delete", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully deleted")
}
