package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventHandler struct {
	Repo *repository.EventRepo
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, h.Repo)
}

func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	GetByID(w, r, h.Repo)
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	if AdminCheck(w, r) {
		return
	}

	// Parse
	var request struct {
		Name        string    `json:"name" bson:"name,omitempty" validate:"required,min=1,max=40"`
		Description string    `json:"description" bson:"description" validate:"required"`
		Time        time.Time `json:"time" bson:"time" validate:"required"`
	}
	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	Create(w, r, h.Repo, &models.Event{
		Name:        request.Name,
		Description: request.Description,
		Time:        request.Time,
	})
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Parse
	var request struct {
		Name        string    `json:"name" bson:"name,omitempty" validate:"omitempty,admin,omitempty,min=1,max=40"`
		Description string    `json:"description" bson:"description,omitempty" validate:"omitempty,admin,omitempty"`
		Time        time.Time `json:"time" bson:"time,omitempty" validate:"omitempty,admin,omitempty"`
	}
	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	// Do work
	err := h.Repo.Update(r.Context(), parsedId, request)
	if utils.CheckError(w, err, "Failed to update", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully updated")
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Check access
	if userAuth.Role == models.Admin {
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
