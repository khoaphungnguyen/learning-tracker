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

type LearningGoal struct {
	ID              int             `json:"id"`
	Title           string          `json:"title"`
	LearningEntries []LearningEntry `json:"learningEntry"`
}

// Placeholder for a database(temporary)
var learningEntries = []LearningEntry{
	{ID: 0, Title: "Learning Binary Search", Description: "O(logn) time", Completed: false},
	{ID: 1, Title: "Learning Linked List", Description: "O(n) time", Completed: false},
}

var learningGoals = []LearningGoal{
	{ID: 0, Title: "Learning Algorithms", LearningEntries: learningEntries},
}

func handleAllLearningEntries(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	goalIDParam := queryValues.Get("goal-id")
	entryIDParam := queryValues.Get("entry-id")
	goalID, err := strconv.Atoi(goalIDParam)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}
	entryID, err := strconv.Atoi(entryIDParam)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	if entryIDParam == "" && r.Method == http.MethodGet {
		// If no 'id' query parameter is provided, fetch and return all entries
		fmt.Println("GET Method: Get all records")
		json.NewEncoder(w).Encode(learningGoals[goalID].LearningEntries)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Retrieve and return a single learning entry
		for _, entry := range learningGoals[goalID].LearningEntries {
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
		for index, entry := range learningGoals[goalID].LearningEntries {
			if entry.ID == entryID {
				learningGoals[goalID].LearningEntries[index].Title = updatedEntry.Title
				learningGoals[goalID].LearningEntries[index].Description = updatedEntry.Description
				fmt.Println("PUT Method: Changed a record")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(fmt.Sprintf("%d", entry.ID)))
				return
			}
		}
		http.Error(w, "Entry not found", http.StatusNotFound)

	case http.MethodDelete:
		// Delete the entry by ID
		for index, entry := range learningGoals[goalID].LearningEntries {
			if entry.ID == entryID {
				learningGoals[goalID].LearningEntries = append(learningGoals[goalID].LearningEntries[:index], learningGoals[goalID].LearningEntries[index+1:]...)
				fmt.Println("DELETE Method: Delete a record")
				w.WriteHeader(http.StatusNoContent)
				w.Write([]byte(fmt.Sprintf("%d", entry.ID)))
				return
			}
			http.Error(w, "Entry not found", http.StatusNotFound)
		}
	}
}

func handleAddNewEntry(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	goalIDParam := queryValues.Get("goal-id")
	goalID, err := strconv.Atoi(goalIDParam)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}
	var newEntry LearningEntry
	err = json.NewDecoder(r.Body).Decode(&newEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(learningGoals[goalID].LearningEntries) == 0 {
		learningGoals[goalID].LearningEntries = append(learningGoals[goalID].LearningEntries, newEntry)
		fmt.Println("POST Method: Add a record")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("New entry added successfully with ID %d", newEntry.ID)))
		return
	}
	newEntry.ID = len(learningGoals[goalID].LearningEntries) + 1
	learningGoals[goalID].LearningEntries = append(learningGoals[goalID].LearningEntries, newEntry)
	fmt.Println("POST Method: Add a record")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New entry added successfully with ID %d", newEntry.ID)))
}

func handleAddNewGoal(w http.ResponseWriter, r *http.Request) {
	var newGoal LearningGoal
	err := json.NewDecoder(r.Body).Decode(&newGoal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(len(learningGoals))
	newGoal.ID = len(learningGoals)
	learningGoals = append(learningGoals, newGoal)
	fmt.Println("POST Method: Add a record")
	w.WriteHeader(http.StatusCreated)
	// return the ID of the newly created entry
	w.Write([]byte(fmt.Sprintf("%d", newGoal.ID)))
}

func handleAllGoal(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	goalIDParam := queryValues.Get("id")

	if goalIDParam == "" && r.Method == http.MethodGet {
		// If no 'id' query parameter is provided, fetch and return all entries
		fmt.Println("GET Method: Get all records")
		json.NewEncoder(w).Encode(learningGoals)
		return
	}

	// If 'id' query parameter is provided, fetch and return the specific entry
	goalID, err := strconv.Atoi(goalIDParam)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Retrieve and return a single learning entry
		for _, goal := range learningGoals {
			if goal.ID == goalID {
				fmt.Println("GET Method: Get one records")
				json.NewEncoder(w).Encode(goal)
				return
			}
		}
		http.Error(w, "Entry not found", http.StatusNotFound)
	case http.MethodPut:
		// Update the entry by ID
		var updatedGoal LearningGoal
		err := json.NewDecoder(r.Body).Decode(&updatedGoal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for index, goal := range learningGoals {
			if goal.ID == goalID {
				learningGoals[index] = updatedGoal
				fmt.Println("PUT Method: Changed a record")
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "Entry not found", http.StatusNotFound)

	case http.MethodDelete:
		// Delete the entry by ID
		for index, goal := range learningGoals {
			if goal.ID == goalID {
				learningGoals = append(learningGoals[:index], learningGoals[index+1:]...)
				fmt.Println("DELETE Method: Delete a record")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.Error(w, "Entry not found", http.StatusNotFound)
		}
	}
}
