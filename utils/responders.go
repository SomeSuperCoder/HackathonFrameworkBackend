package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, value any) {
	resultString, err := json.Marshal(value)
	if CheckError(w, err, "Failed to serialize JSON", http.StatusInternalServerError) {
		return
	}

	fmt.Fprintln(w, string(resultString))
}
