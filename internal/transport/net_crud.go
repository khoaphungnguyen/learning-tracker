package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"crypto/rand"
	"golang.org/x/crypto/argon2"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
)


// Handle new user
func (h *NetHandler) handleNewUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	// Generate Hash
	hash := argon2.IDKey([]byte(user.PasswordHash), salt, 1, 64*1024, 4, 32)

	newId, err := h.netHandler.CreateUser(user.Username,string(hash), user.FirstName, user.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New user is added successfully: EntryID#%d", newId)))
}

func (h *NetHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	// Implementing CORS for frontend compatibility
	w.Header().Set("Access-Control-Allow-Origin", "*")



	switch r.Method {
	// Retrieve and return a single user entry, now based on the authenticated user
	case http.MethodGet:
		// Replacing hardcoded userID with one from request context
		user, err := h.netHandler.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
		return

	// Update the user by ID
	case http.MethodPut:
		var updatedUser models.User
		err := json.NewDecoder(r.Body).Decode(&updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.netHandler.UpdateUser(userID, updatedUser.Username, updatedUser.FirstName, updatedUser.LastName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Updated User#%d successfully", userID)))
		return

	// Delete the user by ID
	case http.MethodDelete:
		err := h.netHandler.DeleteUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}




// Handle new goal
func (h *NetHandler) handleNewGoal(w http.ResponseWriter, r *http.Request) {
	var newGoal models.LearningGoals
	err := json.NewDecoder(r.Body).Decode(&newGoal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newID, err := h.netHandler.CreateGoal(newGoal.UserID, newGoal.Title, newGoal.StartDate, newGoal.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", newID)))
}


// 	UpdateGoal(id int, userID int, title string, startDate time.Time, endDate time.Time) error
// 	DeleteGoal(id int) error
// 	GetAllGoalsByUserID(userID int) ([]models.LearningGoals, error)
// 	GetGoalByID(id int) (models.LearningGoals, error)


// Handle all goals
func (h *NetHandler) handleGoals(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	goalIDParam := queryValues.Get("userID")

	// If no 'id' query parameter is provided, fetch and return all goals
	if r.Method == http.MethodGet && goalIDParam == "" {
		learningGoals, err := h.netHandler.GetAllGoalsByUserID(user.id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(learningGoals)
		return
	}

	// If 'id' query parameter is provided, fetch and return the specific entry
	goalID, err := strconv.Atoi(goalIDParam)
	if err != nil {
		http.Error(w, "Invalid GoalID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// Retrieve and return a single learning goal
	case http.MethodGet:
		goal, err := h.netHandler.GetGoalByID(goalID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(goal)
		return

	// Update the goal by goalID
	case http.MethodPut:
		var updatedGoal models.LearningGoals
		err := json.NewDecoder(r.Body).Decode(&updatedGoal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.netHandler.UpdateGoal(goalID, updatedGoal.Title, updatedGoal.StartDate, updatedGoal.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Updated Goal#%d successfully", goalID)))
		return

	// Delete the entry by ID
	case http.MethodDelete:
		err := h.netHandler.DeleteGoal(goalID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}



// Handle new entry
func (h *NetHandler) handleNewEntry(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	goalID, err := strconv.Atoi(queryValues.Get("goalID"))
	if err != nil {
		http.Error(w, "Invalid goalID", http.StatusBadRequest)
		return
	}
	var newEntry models.LearningEntry
	err = json.NewDecoder(r.Body).Decode(&newEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newId, err := h.netHandler.CreateEntry(goalID, newEntry.Title, newEntry.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New entry is added successfully: EntryID#%d", newId)))
}

// Handle all learning entries
func (h *NetHandler) handleEntries(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	entryIDParam := queryValues.Get("entryID")

	// If no 'id' query parameter is provided, fetch and return all entries
	if r.Method == http.MethodGet && entryIDParam == "" {
		goalID, err := strconv.Atoi(queryValues.Get("goalID"))
		if err != nil {
			http.Error(w, "Invalid GoalID", http.StatusBadRequest)
			return
		}
		entries, err := h.netHandler.GetAllEntriesByGoalID(goalID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(entries)
		
		return
	}

	entryID, err := strconv.Atoi(entryIDParam)
	if err != nil {
		http.Error(w, "Invalid EntryID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	// Retrieve and return a single learning entry
	case http.MethodGet:
		entry, err := h.netHandler.GetEntryByID(entryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(entry)
		return

	// Update the entry by ID
	case http.MethodPut:
		var updatedEntry models.LearningEntry
		err := json.NewDecoder(r.Body).Decode(&updatedEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.netHandler.UpdateEntry(entryID, updatedEntry.Title, updatedEntry.Description, updatedEntry.Date, updatedEntry.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Updated Entry#%d successfully", entryID)))
		return
	// Delete the entry by ID
	case http.MethodDelete:
		err := h.netHandler.DeleteEntry(entryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}


// Handle user file upload
func (h *NetHandler) handleUserFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["files"]

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()

			fileSize := fileHeader.Size
			maxFileSize := int64(25 << 20) // 25 MB
			if fileSize > maxFileSize {
				http.Error(w, "Uploaded file exceeds maximum file size limit (25MB)", http.StatusRequestEntityTooLarge)
				return
			}

			// Generate a timestamp-based unique filename
			currentTime := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
			uniqueFilename := currentTime + "_" + fileHeader.Filename

			// Save the file using the uniqueFilename
			f, err := os.Create("./uploads/" + uniqueFilename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			// Copy the uploaded file to the created file on the filesystem
			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Write([]byte("File uploaded successfully"))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}
