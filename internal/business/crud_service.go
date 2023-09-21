package business

import (
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
)

// User Operations

func (s *LearningService) CreateUser(username string, passwordHash string, firstName string, lastName string) (int64, error) {
	return s.learningStore.CreateUser(username, passwordHash, firstName, lastName)
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

// Learning Entry Operations

func (s *LearningService) GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error) {
	return s.learningStore.GetAllEntriesByGoalID(goalID)
}

func (s *LearningService) GetEntryByID(id int) (models.LearningEntry, error) {
	return s.learningStore.GetEntryByID(id)
}

func (s *LearningService) CreateEntry(goalID int, title string, description string) (int64, error) {
	return s.learningStore.CreateEntry(goalID, title, description)
}

func (s *LearningService) UpdateEntry(id int, title string, description string, date time.Time, status string) error {
	return s.learningStore.UpdateEntry(id, title, description, date, status)
}

func (s *LearningService) DeleteEntry(id int) error {
	return s.learningStore.DeleteEntry(id)
}

// Learning File Operations

func (s *LearningService) GetAllFilesByGoalID(goalID int) ([]models.LearningFiles, error) {
	return s.learningStore.GetAllFilesByGoalID(goalID)
}

func (s *LearningService) GetFileByID(id int) (models.LearningFiles, error) {
	return s.learningStore.GetFileByID(id)
}

func (s *LearningService) CreateFile(goalID int, ownerID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error) {
	return s.learningStore.CreateFile(goalID, ownerID, fileName, fileSize, fileType, filePath)
}

func (s *LearningService) UpdateFile(id int, goalID int, ownerID int, fileName string, fileSize int64, fileType string) error {
	return s.learningStore.UpdateFile(id, goalID, ownerID, fileName, fileSize, fileType)
}

func (s *LearningService) DeleteFile(id int) error {
	return s.learningStore.DeleteFile(id)
}
