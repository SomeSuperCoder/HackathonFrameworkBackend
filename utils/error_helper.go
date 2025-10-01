package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func CheckJSONError(w http.ResponseWriter, err error) bool {
	return CheckError(w, err, "Failed to parse JSON", http.StatusBadRequest)
}

func CheckJSONValidError(w http.ResponseWriter, err error) bool {
	return CheckError(w, err, "JSON validation failed", http.StatusBadRequest)
}

func CheckError(w http.ResponseWriter, err error, message string, code int) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf("%v: %v", message, err.Error()), code)
		return true
	}
	return false
}

func CheckErrorDeadly(err error, message string) {
	if err != nil {
		logrus.Fatalf("%v: %v", message, err.Error())
		os.Exit(1)
	}
}
