package learningbusiness

import (
	"time"

	learningmodel "github.com/khoaphungnguyen/learning-tracker/internal/learning/model"
)

// Learning Goal Operations
func (s *LearningService) CreateGoal(userID int, title string, startDate time.Time, endDate time.Time) (int64, error) {
	return s.learningStore.CreateGoal(userID, title, startDate, endDate)
}

func (s *LearningService) UpdateGoal(id int, userID int, title string, startDate time.Time, endDate time.Time) error {
	return s.learningStore.UpdateGoal(id, userID, title, startDate, endDate)
}

func (s *LearningService) DeleteGoal(id int, userID int) error {
	return s.learningStore.DeleteGoal(id, userID)
}

func (s *LearningService) GetAllGoalsByUserID(userID int) ([]learningmodel.LearningGoals, error) {
	return s.learningStore.GetAllGoalsByUserID(userID)
}

func (s *LearningService) GetGoalByID(id int) (learningmodel.LearningGoals, error) {
	return s.learningStore.GetGoalByID(id)
}

// Learning Entry Operations
func (s *LearningService) CreateEntry(goalID int, userID int, title string, description string) (int64, error) {
	return s.learningStore.CreateEntry(goalID, userID, title, description)
}

func (s *LearningService) UpdateEntry(id int, userID int, title string, description string, status string) error {
	return s.learningStore.UpdateEntry(id, userID, title, description, status)
}

func (s *LearningService) DeleteEntry(id int, userID int) error {
	return s.learningStore.DeleteEntry(id, userID)
}
func (s *LearningService) GetAllEntriesByGoalID(goalID int) ([]learningmodel.LearningEntry, error) {
	return s.learningStore.GetAllEntriesByGoalID(goalID)
}

func (s *LearningService) GetEntryByID(id int) (learningmodel.LearningEntry, error) {
	return s.learningStore.GetEntryByID(id)
}

// Learning File Operations
func (s *LearningService) CreateFile(entryID int, userID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error) {
	return s.learningStore.CreateFile(entryID, userID, fileName, fileSize, fileType, filePath)
}

func (s *LearningService) UpdateFile(id int, userID int, fileName string, fileSize int64, fileType string, filePath string) error {
	return s.learningStore.UpdateFile(id, userID, fileName, fileSize, fileType, filePath)
}

func (s *LearningService) DeleteFile(id int, userID int) error {
	return s.learningStore.DeleteFile(id, userID)
}

func (s *LearningService) GetAllFilesByEntryID(entryID int) ([]learningmodel.LearningFiles, error) {
	return s.learningStore.GetAllFilesByEntryID(entryID)
}

func (s *LearningService) GetFileByID(id int) (learningmodel.LearningFiles, error) {
	return s.learningStore.GetFileByID(id)
}
