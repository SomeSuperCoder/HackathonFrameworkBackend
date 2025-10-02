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
	mux.Handle("/users/", loadUserRoutes(db))
	mux.Handle("/teams/", loadTeamRoutes(db))
	mux.Handle("/cases/", loadCaseRoutes(db))
	mux.Handle("/events/", loadEventRoutes(db))

	return middleware.LoggerMiddleware(mux)
}

func loadCaseRoutes(db *mongo.Database) http.Handler {
	caseMux := http.NewServeMux()
	caseHandler := &handlers.CaseHandler{
		Repo: repository.NewCaseRepo(db),
	}

	caseMux.HandleFunc("GET /", caseHandler.Get)
	caseMux.HandleFunc("GET /{id}", caseHandler.GetByID)
	caseMux.HandleFunc("POST /", middleware.AuthMiddleware(caseHandler.Create, db))
	caseMux.HandleFunc("PATCH /{id}", middleware.AuthMiddleware(caseHandler.Update, db))
	caseMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(caseHandler.Delete, db))

	return http.StripPrefix("/cases", caseMux)
}

func loadEventRoutes(db *mongo.Database) http.Handler {
	eventMux := http.NewServeMux()
	eventHandler := &handlers.EventHandler{
		Repo: repository.NewEventRepo(db),
	}

	eventMux.HandleFunc("GET /", eventHandler.Get)
	eventMux.HandleFunc("GET /{id}", eventHandler.GetByID)
	eventMux.HandleFunc("POST /", middleware.AuthMiddleware(eventHandler.Create, db))
	eventMux.HandleFunc("PATCH /{id}", middleware.AuthMiddleware(eventHandler.Update, db))
	eventMux.HandleFunc("DELETE /{id}", middleware.AuthMiddleware(eventHandler.Delete, db))

	return http.StripPrefix("/events", eventMux)
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

func loadUserRoutes(db *mongo.Database) http.Handler {
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
