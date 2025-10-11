package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
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

	var request struct {
		Name        string `json:"name" bson:"name" validate:"required,min=1,max=40"`
		Description string `json:"description" bson:"description" validate:"required"`
		ImageURI    string `json:"image_uri" bson:"image_uri" validate:"omitempty,url"`
	}
	if DefaultParseAndValidate(w, r, &request) {
		return
	}

	Create(w, r, h.Repo, &models.Case{
		Name:        request.Name,
		Description: request.Description,
		ImageURI:    request.ImageURI,
	})
}

func (h *CaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name" bson:"name,omitempty" validate:"omitempty,admin,min=1,max=40"`
		Description string `json:"description" bson:"description,omitempty" validate:"omitempty,admin"`
		ImageURI    string `json:"image_uri" bson:"image_uri,omitempty" validate:"omitempty,admin,url"`
	}
	Update(w, r, h.Repo, request)
}

func (h *CaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	AdminOnlyDelete(w, r, h.Repo)
}
