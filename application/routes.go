package application

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/handlers"
	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func loadRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	mux.Handle("/users/", loadAuthRoutes(db))

	return middleware.LoggerMiddleware(mux)
}

func loadAuthRoutes(db *mongo.Database) http.Handler {
	authMux := http.NewServeMux()
	authHandler := &handlers.UserHandler{
		Repo: *repository.NewUserRepo(db),
	}

	authMux.HandleFunc("GET /{id}", authHandler.GetUser)
	authMux.HandleFunc("GET /by-name/{username}", authHandler.GetUserByUsername)

	return http.StripPrefix("/users", authMux)
}
