package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"github.com/go-playground/validator/v10"
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

var validate = validator.New()

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
	if utils.CheckError(w, err, "Failed to fetch messages", http.StatusInternalServerError) {
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
	id := r.PathValue("id")

	// Parse
	parsedId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID provided", http.StatusBadRequest)
		return
	}

	// Do work
	team, err := h.Repo.GetByID(r.Context(), parsedId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, fmt.Sprintf("Not found: %s", err.Error()), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get from DB: %s", err.Error()), http.StatusInternalServerError)
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
	id := r.PathValue("id")

	// Parse
	parsedId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID provided", http.StatusBadRequest)
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
		Name string `json:"name" bson:"name,omitempty" validate:"omitempty,min=1,max=20"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	// Validate
	err = validate.Struct(request)
	if utils.CheckJSONValidError(w, err) {
		return
	}

	// Do work
	err = h.Repo.Create(r.Context(), &models.Team{
		Name:            request.Name,
		Leader:          userAuth.ID,
		Repos:           make([]string, 0),
		PresentationURI: "",
		Grades:          make(map[bson.ObjectID]map[bson.ObjectID]int),
		CratedAt:        time.Now(),
	})

	if utils.CheckError(w, err, "Failed to create", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully created")
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	if userAuth.Role == models.Admin || userAuth.Role == models.Judge || parsedId == userAuth.Team {
	} else {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Parse TODO: validate URLs and grades
	var request struct {
		Name            string                                  `json:"name" bson:"name,omitempty" validate:"omitempty,min=1,max=20"`
		Leader          bson.ObjectID                           `json:"leader" bson:"leader,omitempty"`
		Repos           []string                                `json:"repos" bson:"repos,omitempty"`
		PresentationURI string                                  `json:"presentation_uri" bson:"presentation_uri,omitempty"`
		Grades          map[bson.ObjectID]map[bson.ObjectID]int `json:"grades" bson:"grades,omitempty"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	// Validate
	err = validate.Struct(request)
	if utils.CheckJSONValidError(w, err) {
		return
	}

	// Extra access checks TODO: use custom validator tags for that
	if request.Leader != bson.NilObjectID {
		team, err := h.Repo.GetByID(r.Context(), parsedId)
		if utils.CheckError(w, err, "Failed to get team from DB", http.StatusInternalServerError) {
			return
		}
		if team.Leader == userAuth.ID {
		} else {
			http.Error(w, "Only the currernt leader can set the new one", http.StatusForbidden)
			return
		}
	}
	if request.Grades != nil {
		if userAuth.Role == models.Admin || userAuth.Role == models.Judge {
		} else {
			http.Error(w, "You cannot change grades", http.StatusForbidden)
			return
		}
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
	if userAuth.Role == models.Admin || parsedId == userAuth.Team {
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
