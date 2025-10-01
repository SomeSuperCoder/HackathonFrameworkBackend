package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/sirupsen/logrus"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request, db *mongo.Database) (*repository.UserAuth, error, int) {
	// Init repo
	repo := repository.UserRepo{
		Database: db,
	}

	// Load init data from header
	initData := r.Header.Get("TG-Init-Data")

	token := os.Getenv("TELEGRAM_TOKEN")
	expIn := 24 * time.Hour

	err := initdata.Validate(initData, token, expIn)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate initdata: %w", err), http.StatusBadRequest
	}

	// Parse initdata
	initDataParsed, err := initdata.Parse(initData)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse initdata: %w", err), http.StatusBadRequest
	}
	username := initDataParsed.User.Username
	logrus.Info(username)

	user, err := repo.GetUserByUsername(r.Context(), username)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err), http.StatusUnauthorized
	}
	userAuth := &repository.UserAuth{
		Username: user.Username,
		UserID:   user.ID,
	}

	return userAuth, nil, http.StatusNoContent
}
