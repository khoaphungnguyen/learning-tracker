package api

import (
	"encoding/json"
	"fmt"
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

func handleAllLearningEntries(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	entryIDParam := queryValues.Get("id")

	if entryIDParam == "" && r.Method == http.MethodGet {
		// If no 'id' query parameter is provided, fetch and return all entries
		fmt.Println("GET Method: Get all records")
		json.NewEncoder(w).Encode(learningEntries)
		return
	}

	// If 'id' query parameter is provided, fetch and return the specific entry
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
				fmt.Println("GET Method: Get one records")
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
				fmt.Println("PUT Method: Changed a record")
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
				fmt.Println("DELETE Method: Delete a record")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.Error(w, "Entry not found", http.StatusNotFound)
		}
	}
}

func handleAddNewEntry(w http.ResponseWriter, r *http.Request) {
	var newEntry LearningEntry
	err := json.NewDecoder(r.Body).Decode(&newEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newEntry.ID = len(learningEntries) + 1
	learningEntries = append(learningEntries, newEntry)
	fmt.Println("POST Method: Add a record")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New entry added successfully with ID %d", newEntry.ID)))
}
