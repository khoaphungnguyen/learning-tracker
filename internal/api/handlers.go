package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type LearningEntry struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Placeholder for a database(temporary)
var learningEntries = []LearningEntry{
	{ID: 1, Title: "Learning Go", Description: "Learn Go", Completed: false},
	{ID: 2, Title: "Learning Python", Description: "Learn Python", Completed: false},
}

// Handler function for API routes
func handleAllLearningEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Retrieve and return all learning entries
		json.NewEncoder(w).Encode(learningEntries)
	case http.MethodPost:
		// Create a new learning entry
		var newEntry LearningEntry
		err := json.NewDecoder(r.Body).Decode(&newEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		newEntry.ID = len(learningEntries) + 1
		learningEntries = append(learningEntries, newEntry)
		w.WriteHeader(http.StatusCreated)
	}

}

func handleSingleLearningEntry(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	entryIDParam := queryValues.Get("id")
	entryID, err := strconv.Atoi(entryIDParam)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Retrieve and return a single learning entry
		for _, entry := range learningEntries {
			if entry.ID == entryID {
				json.NewEncoder(w).Encode(entry)
				return
			}
		}
		http.Error(w, "Entry not found", http.StatusNotFound)
	case http.MethodPut:
		// Update the entry by ID
		var updatedEntry LearningEntry
		err := json.NewDecoder(r.Body).Decode(&updatedEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for index, entry := range learningEntries {
			if entry.ID == entryID {
				learningEntries[index] = updatedEntry
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "Entry not found", http.StatusNotFound)

	case http.MethodDelete:
		// Delete the entry by ID
		for index, entry := range learningEntries {
			if entry.ID == entryID {
				learningEntries = append(learningEntries[:index], learningEntries[index+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.Error(w, "Entry not found", http.StatusNotFound)
		}
	}
}
