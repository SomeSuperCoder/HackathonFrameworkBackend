package handlers

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CaseHandler struct {
	Repo *repository.CaseRepo
}

func (h *CaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, h.Repo)
}

func (h *CaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	GetByID(w, r, h.Repo)
}

func (h *CaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if AdminCheck(w, r) {
		return
	}

	// Parse
	var request struct {
		Name        string `json:"name" bson:"name" validate:"required,min=1,max=40"`
		Description string `json:"description" bson:"description" validate:"required"`
		ImageURI    string `json:"image_uri" bson:"image_uri" validate:"omitempty,url"`
	}
	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	// Do work
	err := h.Repo.Create(r.Context(), &models.Case{
		Name:        request.Name,
		Description: request.Description,
		ImageURI:    request.ImageURI,
	})

	if utils.CheckError(w, err, "Failed to create", http.StatusInternalServerError) {
		return
	}

	// Respond
	fmt.Fprintf(w, "Successfully created")
}

func (h *CaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Parse
	var request struct {
		Name        string `json:"name" bson:"name,omitempty" validate:"omitempty,admin,omitempty,admin,min=1,max=40"`
		Description string `json:"description" bson:"description,omitempty" validate:"omitempty,admin,omitempty,admin"`
		ImageURI    string `json:"image_uri" bson:"image_uri,omitempty" validate:"omitempty,admin,omitempty,admin,url"`
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

func (h *CaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
