package utils

import (
	"net/http"
)

type Check struct {
	Condition   bool
	Requirement bool
	Message     string
}

func MultiAccessCheck(w http.ResponseWriter, checks []Check) bool {
	for _, check := range checks {
		if AccessCheck(w, check.Condition, check.Requirement, check.Message) {
			return true
		}
	}
	return false
}

func AccessCheck(w http.ResponseWriter, condition bool, requirement bool, message string) bool {
	if condition && !requirement {
		http.Error(w, message, http.StatusForbidden)
		return true
	}
	return false
}
