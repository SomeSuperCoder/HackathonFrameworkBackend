package utils

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ParseRequestID(w http.ResponseWriter, r *http.Request) (bson.ObjectID, bool) {
	id := r.PathValue("id")

	// Parse
	parsedID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID provided", http.StatusBadRequest)
		return bson.NewObjectID(), true
	}

	return parsedID, false
}
