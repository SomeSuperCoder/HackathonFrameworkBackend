package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
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
	var request struct {
		Text string `json:"text" bson:"text,omitempty" validate:"omitempty,admin,required,min=1,max=40"`
	}
	Update(w, r, h.Repo, request)
}

func (h *CriterionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, h.Repo)
}
