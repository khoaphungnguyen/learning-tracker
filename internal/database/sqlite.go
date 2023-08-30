package database

import (
	"database/sql"
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type dbService struct {
	DB *sql.DB
}

func NewDBService() (*dbService, error) {
	db, err := sql.Open("sqlite3", "./learning.db")
	if err != nil {
		panic(err)
	}
	return &dbService{DB: db}, nil
}

func (service *dbService) CreateTable() error {
	// Create the learning_goals table
	_, err := service.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_goals (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			startdate DATETIME,
			enddate DATETIME
		)
	`)
	if err != nil {
		panic(err)
	}

	// Create the learning_entries table
	_, err = service.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_entries (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			goal_id INTEGER,
			title TEXT,
			description TEXT,
			date DATETIME,
			completed BOOLEAN,
			FOREIGN KEY (goal_id) REFERENCES learning_goals(id)
		)
	`)
	if err != nil {
		panic(err)
	}

	return nil
}

func (service *dbService) CreateEntry(goalID int, title string, description string, date time.Time, completed bool) error {
	// Insert the entry
	_, err := service.DB.Exec(`
		INSERT INTO learning_entries (goal_id, title, description, date, completed) VALUES (?, ?, ?, ?, ?)
	`, goalID, title, description, date, completed)

	if err != nil {
		panic(err)
	}
	return nil
}

func (service *dbService) UpdateEntry(id int, title string, description string, date time.Time, completed bool) error {
	// Update the entry
	_, err := service.DB.Exec(`
		UPDATE learning_entries SET title=?, description=?, date=?, completed=? WHERE id=?
	`, title, description, date, completed, id)

	if err != nil {
		panic(err)
	}
	return nil
}

func (service *dbService) DeleteEntry(id int) error {
	// Delete the entry
	_, err := service.DB.Exec(`
		DELETE FROM learning_entries WHERE id=?
	`, id)

	if err != nil {
		panic(err)
	}
	return nil
}

func (service *dbService) GetEntry(id int) (models.LearningEntry, error) {
	// Get the entry
	var entry models.LearningEntry
	err := service.DB.QueryRow(`
		SELECT id, title, description, date, completed FROM learning_entries WHERE id=?
	`, id).Scan(&entry.ID, &entry.Title, &entry.Description, &entry.Date, &entry.Completed)

	if err != nil {
		panic(err)
	}
	return entry, nil
}

func (service *dbService) GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error) {
	// Get all entries for the specified goalID
	var entries []models.LearningEntry
	rows, err := service.DB.Query(`
		SELECT id, title, description, date, completed FROM learning_entries WHERE goal_id=?
	`, goalID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.LearningEntry
		err = rows.Scan(&entry.ID, &entry.Title, &entry.Description, &entry.Date, &entry.Completed)
		if err != nil {
			panic(err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (service *dbService) CreateGoal(title string, startDate time.Time, endDate time.Time) error {
	// Insert the goal
	_, err := service.DB.Exec(`
		INSERT INTO learning_goals (title, startdate, enddate) VALUES (?, ?, ?)
	`, title, startDate, endDate)

	if err != nil {
		panic(err)
	}
	return nil
}

func (service *dbService) UpdateGoal(id int, title string, startDate time.Time, endDate time.Time) error {
	// Update the goal
	_, err := service.DB.Exec(`
		UPDATE learning_goals SET title=?, startdate=?, enddate=? WHERE id=?
	`, title, startDate, endDate, id)

	if err != nil {
		panic(err)
	}
	return nil
}

func (service *dbService) DeleteGoal(id int) error {
	// Delete the goal
	_, err := service.DB.Exec(`
		DELETE FROM learning_goals WHERE id=?
	`, id)

	if err != nil {
		panic(err)
	}
	return nil
}

func (service *dbService) GetGoal(id int) (models.LearningGoals, error) {
	// Get the goal
	var goal models.LearningGoals
	err := service.DB.QueryRow(`
		SELECT id, title, startdate, enddate FROM learning_goals WHERE id=?
	`, id).Scan(&goal.ID, &goal.Title, &goal.StartDate, &goal.EndDate)

	if err != nil {
		panic(err)
	}
	return goal, nil
}

func (service *dbService) GetAllGoals() ([]models.LearningGoals, error) {
	// Get all goals
	var goals []models.LearningGoals
	rows, err := service.DB.Query(`
		SELECT id, title, startdate, enddate FROM learning_goals
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var goal models.LearningGoals
		err = rows.Scan(&goal.ID, &goal.Title, &goal.StartDate, &goal.EndDate)
		if err != nil {
			panic(err)
		}
		goals = append(goals, goal)
	}

	return goals, nil
}
