package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserHandler struct {
	Repo *repository.UserRepo
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	// Parse the user id
	parsedUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID provided", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByID(r.Context(), parsedUserID)
	getCommon(user, err, w)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	user, err := h.Repo.GetByUsername(r.Context(), username)
	getCommon(user, err, w)
}

func getCommon(user *models.User, err error, w http.ResponseWriter) {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, fmt.Sprintf("Not found: %s", err.Error()), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get from DB: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	serializedUser, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(serializedUser))
}
