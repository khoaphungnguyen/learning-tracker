package learningstorage

import (
	"database/sql"
	"time"

	learningmodel "github.com/khoaphungnguyen/learning-tracker/internal/learning/model"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type LearningStore interface {

	// Learning goal operations
	CreateGoal(userID int, title string, startDate time.Time, endDate time.Time) (int64, error)
	UpdateGoal(id int, userID int, title string, startDate time.Time, endDate time.Time) error
	DeleteGoal(id int, userID int) error
	GetAllGoalsByUserID(userID int) ([]learningmodel.LearningGoals, error)
	GetGoalByID(id int) (learningmodel.LearningGoals, error)

	// Learning entry operations
	CreateEntry(goalID int, user_id int, title string, description string) (int64, error)
	UpdateEntry(id int, userID int, title string, description string, status string) error
	DeleteEntry(id int, userID int) error
	GetAllEntriesByGoalID(goalID int) ([]learningmodel.LearningEntry, error)
	GetEntryByID(id int) (learningmodel.LearningEntry, error)

	// Learning file operations
	CreateFile(entryID int, userID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error)
	UpdateFile(id int, userID int, fileName string, fileSize int64, fileType string, filePath string) error
	DeleteFile(id int, userID int) error
	GetAllFilesByEntryID(entryID int) ([]learningmodel.LearningFiles, error)
	GetFileByID(id int) (learningmodel.LearningFiles, error)
}

type learningStore struct {
	DB *sql.DB
}

func NewLearningStore() (*learningStore, error) {
	DB, err := sql.Open("sqlite3", "./migrations/learning.db")
	if err != nil {
		panic(err)
	}
	return &learningStore{DB: DB}, nil
}
