package transport

import (
	"net/http"
)

// SetupRoutes sets up all the routes for the API
func (h *NetHandler) SetupRoutes(router *http.ServeMux) {
	// Routes for users
	router.HandleFunc("/api/users", h.handleUsers)
	router.HandleFunc("/api/users/add", h.handleNewUser)

	// Routes for entries
	router.HandleFunc("/api/entries", h.handleEntries)
	router.HandleFunc("/api/entries/add", h.handleNewEntry)

	// Routes for goals
	router.HandleFunc("/api/goals", h.handleGoals)
	router.HandleFunc("/api/goals/add", h.handleNewGoal)

	// Routes for file operations
	router.HandleFunc("/api/files", h.handleFiles)
	router.HandleFunc("/api/files/add", h.handleNewFile)
}
