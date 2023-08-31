package storage

import (
	"database/sql"
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type LearningStore interface {
	CreateTable() error
	CreateEntry(goalID int, title string, description string, date time.Time, completed bool) (int64, error)
	CreateGoal(title string, startDate time.Time, endDate time.Time) (int64, error)
	UpdateEntry(id int, title string, description string, date time.Time, completed bool) error
	UpdateGoal(id int, title string, startDate time.Time, endDate time.Time) error
	DeleteEntry(id int) error
	DeleteGoal(id int) error
	GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error)
	GetAllGoals() ([]models.LearningGoals, error)
	GetGoalByID(id int) (models.LearningGoals, error)
	GetEntryByID(id int) (models.LearningEntry, error)
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
