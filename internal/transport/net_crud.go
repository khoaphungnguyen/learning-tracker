package transport

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/dgrijalva/jwt-go"
	"github.com/khoaphungnguyen/learning-tracker/internal/models"
)

type MyClaims struct {
	Fullname string `json:"fullname"`
	jwt.StandardClaims
}

// checkAuthorization checks if the user is authorized to perform the operation
func checkAuthorization(r *http.Request) (int, bool) {
	userIDFromContext := r.Context().Value("userID").(string)
	queryValues := r.URL.Query()
	userIDParam := queryValues.Get("userID")
	if userIDFromContext != userIDParam {
		return 0, false
	}
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return 0, false
	}
	return userID, true
}

// Handle User's Signin
func (h *NetHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "In valid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve a user by their user name
	user, err := h.netHandler.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	// Verify password hash
	hash := argon2.IDKey([]byte(req.Password), user.Salt, 1, 64*1024, 4, 32)
	if user.Password != string(hash) {
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}
	// Generate JWT token for the valid user
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &MyClaims{
		Fullname: user.LastName,
		StandardClaims: jwt.StandardClaims{
			Audience:  strconv.Itoa(user.ID),
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.JWTKey))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	}

	response := map[string]string{
		"token": tokenString,
	}

	json.NewEncoder(w).Encode(response)

}

// handleUsers manages ADD operations for the User model
func (h *NetHandler) handleAddUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
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
	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)

	_, err = h.netHandler.CreateUser(strings.ToLower(user.Username), hash, salt, user.FirstName, user.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New user is added successfully: %s", user.Username)))
}

// handleUsers manages CRUD operations for the User model
func (h *NetHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authorized to perform the operation
	userID, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodGet:
		user, err := h.netHandler.GetUserByID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
		return

	case http.MethodPut:
		// Assuming the user ID is passed in the URL for update
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
	// Check if the user is authorized to perform the operation
	userID, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var newGoal models.LearningGoals
	err := json.NewDecoder(r.Body).Decode(&newGoal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newID, err := h.netHandler.CreateGoal(userID, newGoal.Title, newGoal.StartDate, newGoal.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", newID)))
}

// handleGoals handles all goal-related operations like GET, POST (for creating new), PUT, DELETE
func (h *NetHandler) handleGoals(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authorized to perform the operation
	userID, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	queryValues := r.URL.Query()
	goalIDParam := queryValues.Get("goalID")

	// If no 'goalID' query parameter is provided, fetch all goals by userID
	if r.Method == http.MethodGet && goalIDParam == "" {
		goals, err := h.netHandler.GetAllGoalsByUserID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(goals)
		return
	}

	goalID, err := strconv.Atoi(goalIDParam)
	if err != nil {
		http.Error(w, "Invalid GoalID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		goal, err := h.netHandler.GetGoalByID(goalID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(goal)
		return

	case http.MethodPost:
		var newGoal models.LearningGoals
		err := json.NewDecoder(r.Body).Decode(&newGoal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newID, err := h.netHandler.CreateGoal(userID, newGoal.Title, newGoal.StartDate, newGoal.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("New goal created with ID#%d", newID)))
		return

	case http.MethodPut:
		var updatedGoal models.LearningGoals
		err := json.NewDecoder(r.Body).Decode(&updatedGoal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.netHandler.UpdateGoal(goalID, updatedGoal.UserID, updatedGoal.Title, updatedGoal.StartDate, updatedGoal.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Updated Goal#%d successfully", goalID)))
		return

	case http.MethodDelete:
		err := h.netHandler.DeleteGoal(goalID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Handle new entry
func (h *NetHandler) handleNewEntry(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authorized to perform the operation
	_, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var newEntryRequest models.LearningEntry
	err := json.NewDecoder(r.Body).Decode(&newEntryRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newId, err := h.netHandler.CreateEntry(newEntryRequest.GoalID, newEntryRequest.UserID, newEntryRequest.Title, newEntryRequest.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New entryID#%d is added successfully", newId)))
}

// handleEntries manages CRUD operations for the LearningEntry model
func (h *NetHandler) handleEntries(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authorized to perform the operation
	userID, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	queryValues := r.URL.Query()
	entryIDParam := queryValues.Get("entryID")

	switch r.Method {
	case http.MethodGet:
		// If no 'id' query parameter is provided, get all entries for a specific goal
		if entryIDParam == "" {
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

		// Get specific entry by ID
		entryID, err := strconv.Atoi(entryIDParam)
		if err != nil {
			http.Error(w, "Invalid EntryID", http.StatusBadRequest)
			return
		}
		entry, err := h.netHandler.GetEntryByID(entryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(entry)
		return

	case http.MethodPost:
		var newEntry models.LearningEntry
		err := json.NewDecoder(r.Body).Decode(&newEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newID, err := h.netHandler.CreateEntry(newEntry.GoalID, newEntry.UserID, newEntry.Title, newEntry.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("New entry created with ID#%d", newID)))
		return

	case http.MethodPut:
		// Assuming the entry ID is passed in the URL for update
		entryID, err := strconv.Atoi(entryIDParam)
		if err != nil {
			http.Error(w, "Invalid EntryID", http.StatusBadRequest)
			return
		}

		var updatedEntry models.LearningEntry
		err = json.NewDecoder(r.Body).Decode(&updatedEntry)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.netHandler.UpdateEntry(entryID, updatedEntry.GoalID, updatedEntry.UserID, updatedEntry.Title, updatedEntry.Description, updatedEntry.Date, updatedEntry.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Updated Entry#%d successfully", entryID)))
		return

	case http.MethodDelete:
		entryID, err := strconv.Atoi(entryIDParam)
		if err != nil {
			http.Error(w, "Invalid EntryID", http.StatusBadRequest)
			return
		}

		err = h.netHandler.DeleteEntry(entryID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// handleNewFile handles file creation
func (h *NetHandler) handleNewFile(w http.ResponseWriter, r *http.Request) {
	_, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var newFile models.LearningFiles
	err := json.NewDecoder(r.Body).Decode(&newFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newID, err := h.netHandler.CreateFile(newFile.EntryID, newFile.UserID, newFile.FileName, newFile.FileSize, newFile.FileType, newFile.FilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New file created with ID#%d", newID)))
}

// handleFiles handles all file-related operations like GET, PUT, DELETE
func (h *NetHandler) handleFiles(w http.ResponseWriter, r *http.Request) {
	_, isAuthorized := checkAuthorization(r)
	if !isAuthorized {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	queryValues := r.URL.Query()
	fileIDParam := queryValues.Get("fileID")

	if r.Method == http.MethodGet && fileIDParam == "" {
		entryID, err := strconv.Atoi(queryValues.Get("entryID"))
		if err != nil {
			http.Error(w, "Invalid EntryID", http.StatusBadRequest)
			return
		}
		files, err := h.netHandler.GetAllFilesByEntryID(entryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(files)
		return
	}

	fileID, err := strconv.Atoi(fileIDParam)
	if err != nil {
		http.Error(w, "Invalid FileID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		file, err := h.netHandler.GetFileByID(fileID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(file)

	case http.MethodPut:
		var updatedFile models.LearningFiles
		err := json.NewDecoder(r.Body).Decode(&updatedFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.netHandler.UpdateFile(fileID, updatedFile.EntryID, updatedFile.UserID,
			updatedFile.FileName, updatedFile.FileSize, updatedFile.FileType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Updated File#%d successfully", fileID)))
		return

	case http.MethodDelete:
		err := h.netHandler.DeleteFile(fileID)
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
