package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const UserKey = "user"

func ExtractUserAuth(r *http.Request) *models.User {
	userAuth, ok := r.Context().Value(UserKey).(*models.User)
	if !ok {
		panic("Failed to extract user auth data")
	}

	return userAuth
}

func AuthMiddleware(next http.HandlerFunc, db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		if os.Getenv("API_TEST") == "" {
			var err error
			var code int
			user, err, code = utils.Authorize(r, db)
			if err != nil {
				http.Error(w, fmt.Errorf("Failed to authorize: %w", err).Error(), code)
				return
			}
		} else {
			user = &models.User{
				ID:        bson.NilObjectID,
				Username:  "test",
				Name:      "Mr. Test",
				Birthdate: time.Now(),
				Role:      models.Admin,
				ChatID:    0,
				Team:      bson.NilObjectID,
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
