package handlers

import (
	"net/http"
	"time"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
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
	// Parse
	var request struct {
		Name        string    `json:"name" bson:"name,omitempty" validate:"omitempty,admin,omitempty,min=1,max=40"`
		Description string    `json:"description" bson:"description,omitempty" validate:"omitempty,admin,omitempty"`
		Time        time.Time `json:"time" bson:"time,omitempty" validate:"omitempty,admin,omitempty"`
	}
	Update(w, r, h.Repo, request)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, h.Repo)
}
