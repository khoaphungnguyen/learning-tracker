package api

import (
	"net/http"
)

// SetupRoutes sets up all the routes for the API
func SetupRoutes(router *http.ServeMux) {
	//Define routes
	router.HandleFunc("/api/learning_entries", handleAllLearningEntries)
	router.HandleFunc("/api/learning_entries/{id}", handleSingleLearningEntry)

}
