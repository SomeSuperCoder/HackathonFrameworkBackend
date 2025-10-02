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
	mux.Handle("/teams/", loadTeamRoutes(db))

	return middleware.LoggerMiddleware(mux)
}

func loadTeamRoutes(db *mongo.Database) http.Handler {
	teamMux := http.NewServeMux()
	teamHandler := &handlers.TeamHandler{
		Repo: repository.NewTeamRepo(db),
	}

	teamMux.HandleFunc("GET /", teamHandler.GetPaged)
	teamMux.HandleFunc("GET /{id}", teamHandler.GetByID)
	teamMux.HandleFunc("GET /{id}/members", teamHandler.GetMembers)
	teamMux.HandleFunc("POST /", middleware.AuthMiddleware(teamHandler.Create, db))
	teamMux.HandleFunc("PATCH /{id}", middleware.AuthMiddleware(teamHandler.Update, db))
	teamMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(teamHandler.Delete, db))

	return http.StripPrefix("/teams", teamMux)
}

func loadAuthRoutes(db *mongo.Database) http.Handler {
	userMux := http.NewServeMux()
	userHandler := &handlers.UserHandler{
		Repo: repository.NewUserRepo(db),
	}

	userMux.HandleFunc("GET /", userHandler.GetPaged)
	userMux.HandleFunc("GET /{id}", userHandler.GetByID)
	userMux.HandleFunc("GET /by-name/{username}", userHandler.GetByUsername)
	userMux.HandleFunc("PATCH /{id}", middleware.AuthMiddleware(userHandler.Update, db))
	userMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(userHandler.Delete, db))

	return http.StripPrefix("/users", userMux)
}
