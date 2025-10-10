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

type CriterionHandler struct {
	Repo *repository.CriterionRepo
}

func (h *CriterionHandler) Get(w http.ResponseWriter, r *http.Request) {
	Get(w, r, h.Repo)

}

func (h *CriterionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	GetByID(w, r, h.Repo)
}

func (h *CriterionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if AdminCheck(w, r) {
		return
	}

	var request struct {
		Text string `json:"text" bson:"text" validate:"required,min=1,max=40"`
	}
	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	Create(w, r, h.Repo, &models.Criterion{
		Text: request.Text,
	})
}

func (h *CriterionHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Load data
	var parsedId bson.ObjectID
	var exit bool
	if parsedId, exit = utils.ParseRequestID(w, r); exit {
		return
	}

	// Parse
	var request struct {
		Text string `json:"text" bson:"text,omitempty" validate:"omitempty,admin,required,min=1,max=40"`
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

func (h *CriterionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
