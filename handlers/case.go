package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/internal/validators"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CaseHandler struct {
	Repo *repository.CaseRepo
}

func (h *CaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Do work
	cases, err := h.Repo.Find(r.Context())
	if utils.CheckError(w, err, "Failed to get from DB", http.StatusInternalServerError) {
		return
	}

	// Respond
	resultString, err := json.Marshal(&cases)
	if utils.CheckError(w, err, "Failed to serialize JSON", http.StatusInternalServerError) {
		return
	}

	fmt.Fprintln(w, string(resultString))
}

func (h *CaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	case_, err := h.Repo.GetByID(r.Context(), parsedId)
	if utils.CheckGetFromDB(w, err) {
		return
	}

	serialized, err := json.Marshal(&case_)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(serialized))
}

func (h *CaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Check access
	if userAuth.Role != models.Admin {
		http.Error(w, "Access denied: only the admin can perform this operation", http.StatusForbidden)
		return
	}

	// Parse
	var request struct {
		Name        string `json:"name" bson:"name" validate:"required,min=1,max=40"`
		Description string `json:"description" bson:"description" validate:"required"`
		ImageURI    string `json:"image_uri" bson:"image_uri" validate:"omitempty,url"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	// Validate
	tv := validators.NewAccessValidator(userAuth)
	err = tv.ValidateRequest(&request)
	if utils.CheckJSONValidError(w, err) {
		return
	}

	// Do work
	err = h.Repo.Create(r.Context(), &models.Case{
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

	// Get auth data
	userAuth := middleware.ExtractUserAuth(r)

	// Parse
	var request struct {
		Name        string `json:"name" bson:"name,omitempty" validate:"omitempty,admin,omitempty,admin,min=1,max=40"`
		Description string `json:"description" bson:"description,omitempty" validate:"omitempty,admin,omitempty,admin"`
		ImageURI    string `json:"image_uri" bson:"image_uri,omitempty" validate:"omitempty,admin,omitempty,admin,url"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if utils.CheckJSONError(w, err) {
		return
	}

	tv := validators.NewAccessValidator(userAuth)
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
