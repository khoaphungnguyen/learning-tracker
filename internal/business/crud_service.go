package business

import (
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
)

// User Operations

func (s *LearningService) CreateUser(username string, password []byte, salt []byte, firstName string, lastName string) (int64, error) {
	return s.learningStore.CreateUser(username, password, salt, firstName, lastName)
}

func (s *LearningService) UpdateUser(id int, username string, firstName string, lastName string) error {
	return s.learningStore.UpdateUser(id, username, firstName, lastName)
}

func (s *LearningService) DeleteUser(id int) error {
	return s.learningStore.DeleteUser(id)
}

func (s *LearningService) GetUserByID(id int) (models.User, error) {
	return s.learningStore.GetUserByID(id)
}

func (s *LearningService) GetUserByUsername(username string) (models.User, error) {
	return s.learningStore.GetUserByUsername(username)
}

// Learning Goal Operations
func (s *LearningService) GetGoalByID(id int) (models.LearningGoals, error) {
	return s.learningStore.GetGoalByID(id)
}

func (s *LearningService) CreateGoal(userID int, title string, startDate time.Time, endDate time.Time) (int64, error) {
	return s.learningStore.CreateGoal(userID, title, startDate, endDate)
}

func (s *LearningService) UpdateGoal(id int, userID int, title string, startDate time.Time, endDate time.Time) error {
	return s.learningStore.UpdateGoal(id, userID, title, startDate, endDate)
}

func (s *LearningService) DeleteGoal(id int) error {
	return s.learningStore.DeleteGoal(id)
}

func (s *LearningService) GetAllGoalsByUserID(userID int) ([]models.LearningGoals, error) {
	return s.learningStore.GetAllGoalsByUserID(userID)
}

// Learning Entry Operations

func (s *LearningService) GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error) {
	return s.learningStore.GetAllEntriesByGoalID(goalID)
}

func (s *LearningService) GetEntryByID(id int) (models.LearningEntry, error) {
	return s.learningStore.GetEntryByID(id)
}

func (s *LearningService) CreateEntry(goalID int, userID int, title string, description string) (int64, error) {
	return s.learningStore.CreateEntry(goalID, userID, title, description)
}

func (s *LearningService) UpdateEntry(id int, goalID int, userID int, title string, description string, date time.Time, status string) error {
	return s.learningStore.UpdateEntry(id, goalID, userID, title, description, date, status)
}

func (s *LearningService) DeleteEntry(id int, userID int) error {
	return s.learningStore.DeleteEntry(id, userID)
}

// Learning File Operations

func (s *LearningService) GetAllFilesByEntryID(entryID int) ([]models.LearningFiles, error) {
	return s.learningStore.GetAllFilesByEntryID(entryID)
}

func (s *LearningService) GetFileByID(id int) (models.LearningFiles, error) {
	return s.learningStore.GetFileByID(id)
}

func (s *LearningService) CreateFile(entryID int, userID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error) {
	return s.learningStore.CreateFile(entryID, userID, fileName, fileSize, fileType, filePath)
}

func (s *LearningService) UpdateFile(id int, entryID int, userID int, fileName string, fileSize int64, fileType string) error {
	return s.learningStore.UpdateFile(id, entryID, userID, fileName, fileSize, fileType)
}

func (s *LearningService) DeleteFile(id int) error {
	return s.learningStore.DeleteFile(id)
}
