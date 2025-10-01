package utils

import (
	"fmt"
	"net/http"
)

func CheckJSONError(w http.ResponseWriter, err error) bool {
	return CheckError(w, err, "Failed to parse JSON", http.StatusBadRequest)
}

func CheckError(w http.ResponseWriter, err error, message string, code int) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", message, err.Error()), code)
		return true
	}
	return false
}
