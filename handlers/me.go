package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/global-chat/internal/middleware"
	"github.com/SomeSuperCoder/global-chat/utils"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.ExtractUserAuth(r)

	serialized, err := json.Marshal(user)
	if utils.CheckJSONError(w, err) {
		return
	}

	fmt.Fprintln(w, string(serialized))
}
