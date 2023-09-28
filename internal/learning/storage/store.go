package storage

import (
	"database/sql"
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type LearningStore interface {
	
	// Learning goal operations
	CreateGoal(userID int, title string, startDate time.Time, endDate time.Time) (int64, error)
	UpdateGoal(id int, userID int, title string, startDate time.Time, endDate time.Time) error
	DeleteGoal(id int) error
	GetAllGoalsByUserID(userID int) ([]models.LearningGoals, error)
	GetGoalByID(id int) (models.LearningGoals, error)

	// Learning entry operations
	CreateEntry(goalID int, user_id int, title string, description string) (int64, error)
	UpdateEntry(id int, goalID int , userID int, title string, description string, date time.Time, status string) error
	DeleteEntry(id int, userID int) error
	GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error)
	GetEntryByID(id int) (models.LearningEntry, error)

	// Learning file operations
	CreateFile(entryID int, userID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error)
	UpdateFile(id int, entryID int, userID int, fileName string, fileSize int64, fileType string) error
	DeleteFile(id int) error
	GetAllFilesByEntryID(entryID int) ([]models.LearningFiles, error)
	GetFileByID(id int) (models.LearningFiles, error)
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
