package handlers

import (
	"net/http"
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
	FindPaged(w, r, h.Repo, func(values []models.User, totalCount int64) any {
		return UsersResponse{
			Users:      values,
			TotalCount: totalCount,
		}
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	GetByID(w, r, h.Repo)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	user, err := h.Repo.GetByUsername(r.Context(), username)
	if utils.CheckGetFromDB(w, err) {
		return
	}

	utils.RespondWithJSON(w, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
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
		Name      string          `json:"name" bson:"name,omitempty" validate:"omitempty,self,min=1,max=40"`
		Birthdate time.Time       `json:"birthdate" bson:"birthdate,omitempty" validate:"omitempty,self"`
		Role      models.UserRole `json:"role" bson:"role,omitempty" validate:"omitempty,admin,oneof=0 1 2"`
		Team      bson.ObjectID   `json:"team" bson:"team,omitempty" validate:"omitempty,self"`
	}

	if ParseAndValidate(w, r, validators.NewUserValidator(userAuth, parsedId), &request) {
		return
	}

	UpdateInner(w, r, h.Repo, parsedId, request)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	Delete(w, r, h.Repo, func(w http.ResponseWriter, r *http.Request, id bson.ObjectID, userAuth *models.User) bool {
		if userAuth.Role == models.Admin || id == userAuth.ID {
			return false
		} else {
			http.Error(w, "Access denied", http.StatusForbidden)
			return true
		}
	})
}
