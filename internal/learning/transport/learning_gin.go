package learningtransport

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	learningmodel "github.com/khoaphungnguyen/learning-tracker/internal/learning/model"
)

// Handle new goal from gin
func (h *LearningHandler) CreateGoal(c *gin.Context) {
	// Get user ID from the JWT token
	userID := c.GetInt("id")
	var newGoal learningmodel.LearningGoals
	err := c.BindJSON(&newGoal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newID, err := h.learningHandler.CreateGoal(userID, newGoal.Title, newGoal.StartDate, newGoal.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Goal#%d is created successfully", newID),
	})
}

// Handle update goal from gin
func (h *LearningHandler) UpdateGoal(c *gin.Context) {
	userID := c.GetInt("id")
	var goal learningmodel.LearningGoals
	err := c.BindJSON(&goal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = h.learningHandler.UpdateGoal(goal.ID, userID, goal.Title, goal.StartDate, goal.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Goal#%d is updated successfully", goal.ID),
	})
}

// Handle delete goal from gin
func (h *LearningHandler) DeleteGoal(c *gin.Context) {
	userID := c.GetInt("id")
	goalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid goal ID",
		})
		return
	}
	err = h.learningHandler.DeleteGoal(goalID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Goal#%d is deleted successfully", goalID),
	})
}

// Handle get all goals from gin
func (h *LearningHandler) GetAllGoalsByUserID(c *gin.Context) {
	userID := c.GetInt("id")
	goals, err := h.learningHandler.GetAllGoalsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, goals)
}

// Handle get goal by ID from gin
func (h *LearningHandler) GetGoalByID(c *gin.Context) {
	userID := c.GetInt("id")
	goal, err := h.learningHandler.GetGoalByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, goal)
}

// Handle new entry
func (h *LearningHandler) CreateEntry(c *gin.Context) {
	userID := c.GetInt("id")
	var newEntry learningmodel.LearningEntry
	err := c.BindJSON(&newEntry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newID, err := h.learningHandler.CreateEntry(newEntry.GoalID, userID, newEntry.Title, newEntry.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("New entry#%d is created successfully", newID),
	})
}

// Handle update entry
func (h *LearningHandler) UpdateEntry(c *gin.Context) {
	userID := c.GetInt("id")
	var updatedEntry learningmodel.LearningEntry
	err := c.BindJSON(&updatedEntry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	err = h.learningHandler.UpdateEntry(updatedEntry.ID, userID, updatedEntry.Title, updatedEntry.Description, updatedEntry.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Entry#%d is updated successfully", updatedEntry.ID),
	})
}

func (h *LearningHandler) DeleteEntry(c *gin.Context) {
	userID := c.GetInt("id")
	entryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid entry ID",
		})
		return
	}
	err = h.learningHandler.DeleteEntry(entryID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Entry#%d is deleted successfully", entryID),
	})
}

func (h *LearningHandler) GetAllEntriesByGoalID(c *gin.Context) {
	goalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid goal ID",
		})
		return
	}
	entries, err := h.learningHandler.GetAllEntriesByGoalID(goalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, entries)
}

func (h *LearningHandler) GetEntryByID(c *gin.Context) {
	entryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid entry ID",
		})
		return
	}
	entry, err := h.learningHandler.GetEntryByID(entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, entry)
}

// Handle new file upload and create file
func (h *LearningHandler) CreateFile(c *gin.Context) {
	userID := c.GetInt("id")
	entryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid entry ID",
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileName := file.Filename
	fileSize := file.Size
	fileType := file.Header.Get("Content-Type")
	filePath := fmt.Sprintf("uploads/%d/%s", userID, fileName)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	newID, err := h.learningHandler.CreateFile(entryID, userID, fileName, fileSize, fileType, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("New file#%d is created successfully", newID),
	})
}

// Handle update file and replace the old file
func (h *LearningHandler) UpdateFile(c *gin.Context) {
	userID := c.GetInt("id")
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileName := file.Filename
	fileSize := file.Size
	fileType := file.Header.Get("Content-Type")
	filePath := fmt.Sprintf("uploads/%d/%s", userID, fileName)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	err = h.learningHandler.UpdateFile(fileID, userID, fileName, fileSize, fileType, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File#%d is updated successfully", fileID),
	})
}

// Handle delete file
func (h *LearningHandler) DeleteFile(c *gin.Context) {
	userID := c.GetInt("id")
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}
	file, err := h.learningHandler.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	err = os.Remove(file.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	err = h.learningHandler.DeleteFile(fileID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File#%d is deleted successfully", fileID),
	})
}

// Handle get all files by entry ID
func (h *LearningHandler) GetAllFilesByEntryID(c *gin.Context) {
	entryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid entry ID",
		})
		return
	}
	files, err := h.learningHandler.GetAllFilesByEntryID(entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, files)
}

// Handle get file by ID
func (h *LearningHandler) GetFileByID(c *gin.Context) {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}
	file, err := h.learningHandler.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, file)
}

// Handle download file by ID
func (h *LearningHandler) DownloadFile(c *gin.Context) {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}
	file, err := h.learningHandler.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.File(file.FilePath)
}
