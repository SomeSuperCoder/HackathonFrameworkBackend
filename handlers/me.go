package handlers

import (
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/utils"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, middleware.ExtractUserAuth(r))
}
