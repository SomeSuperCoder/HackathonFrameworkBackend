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
	teamsHandler := &handlers.TeamHandler{
		Repo: repository.NewTeamRepo(db),
	}

	teamMux.HandleFunc("GET /", teamsHandler.GetPaged)
	teamMux.HandleFunc("GET /{id}", teamsHandler.GetByID)
	teamMux.HandleFunc("GET /{id}/members", teamsHandler.GetMembers)
	teamMux.HandleFunc("POST /", middleware.AuthMiddleware(teamsHandler.Create, db))
	teamMux.HandleFunc("PATCH /{id}", middleware.AuthMiddleware(teamsHandler.Update, db))
	teamMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(teamsHandler.Delete, db))

	return http.StripPrefix("/teams", teamMux)
}

func loadAuthRoutes(db *mongo.Database) http.Handler {
	userMux := http.NewServeMux()
	usersHandler := &handlers.UserHandler{
		Repo: repository.NewUserRepo(db),
	}

	userMux.HandleFunc("GET /", usersHandler.GetPaged)
	userMux.HandleFunc("GET /{id}", usersHandler.GetByID)
	userMux.HandleFunc("GET /by-name/{username}", usersHandler.GetByUsername)
	userMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(usersHandler.Delete, db))

	return http.StripPrefix("/users", userMux)
}
