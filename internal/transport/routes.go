package transport

import (
	"net/http"
)

// SetupRoutes sets up all the routes for the API
func (h *NetHandler) SetupRoutes(router *http.ServeMux) {
	// Create a new storage instance

	//Route for entries
	router.HandleFunc("/api/entries", h.handleAllLearningEntries)
	router.HandleFunc("/api/entries/add", h.handleAddNewEntry)

	//Route for goals
	router.HandleFunc("/api/goals", h.handleAllGoal)
	router.HandleFunc("/api/goals/add", h.handleAddNewGoal)
}
