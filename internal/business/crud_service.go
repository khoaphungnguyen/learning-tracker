package business

import (
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
)

func (s *LearningService) GetAllGoals() ([]models.LearningGoals, error) {
	return s.learningStore.GetAllGoals()
}

func (s *LearningService) GetGoalByID(id int) (models.LearningGoals, error) {
	return s.learningStore.GetGoalByID(id)
}

func (s *LearningService) CreateGoal(title string, startDate time.Time, endDate time.Time) (int64, error) {
	return s.learningStore.CreateGoal(title, startDate, endDate)
}

func (s *LearningService) UpdateGoal(id int, title string, startDate time.Time, endDate time.Time) error {
	return s.learningStore.UpdateGoal(id, title, startDate, endDate)
}

func (s *LearningService) DeleteGoal(id int) error {
	return s.learningStore.DeleteGoal(id)
}

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
