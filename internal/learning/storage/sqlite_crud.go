package storage

import (
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)



// CreateGoal inserts a new learning goal into the database and returns its ID.
func (service *learningStore) CreateGoal(userID int, title string, startDate time.Time, endDate time.Time) (int64, error) {
	result, err := service.DB.Exec(`
		INSERT INTO learning_goals (user_id, title, startdate, enddate) VALUES (?, ?, ?, ?)
	`, userID, title, startDate, endDate)
	if err != nil {
		return 0, err
	}
	newID, _ := result.LastInsertId()
	return newID, nil
}

// UpdateGoal updates an existing learning goal.
func (service *learningStore) UpdateGoal(id int, userID int, title string, startDate time.Time, endDate time.Time) error {
	_, err := service.DB.Exec(`
		UPDATE learning_goals SET user_id=?, title=?, startdate=?, enddate=? WHERE id=?
	`, userID, title, startDate, endDate, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGoal deletes a learning goal by ID.
func (service *learningStore) DeleteGoal(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM learning_goals WHERE id=?
	`, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllGoalsByUserID returns all learning goals for a given user ID.
func (service *learningStore) GetAllGoalsByUserID(userID int) ([]models.LearningGoals, error) {
	var goals []models.LearningGoals
	rows, err := service.DB.Query(`
		SELECT id, user_id, title, startdate, enddate FROM learning_goals WHERE user_id=?
	`, userID)
	if err != nil {
		return goals, err
	}
	defer rows.Close()

	for rows.Next() {
		var goal models.LearningGoals
		err = rows.Scan(&goal.ID, &goal.UserID, &goal.Title, &goal.StartDate, &goal.EndDate)
		if err != nil {
			return goals, err
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

// GetGoalByID returns a learning goal by ID.
func (service *learningStore) GetGoalByID(id int) (models.LearningGoals, error) {
	var goal models.LearningGoals
	err := service.DB.QueryRow(`
		SELECT id, user_id, title, startdate, enddate FROM learning_goals WHERE id=?
	`, id).Scan(&goal.ID, &goal.UserID, &goal.Title, &goal.StartDate, &goal.EndDate)
	if err != nil {
		return goal, err
	}
	return goal, nil
}

// CreateEntry inserts a new learning entry into the database and returns its ID.
func (service *learningStore) CreateEntry(goalID int, userID int, title string, description string) (int64, error) {
	result, err := service.DB.Exec(`
        INSERT INTO learning_entries (goal_id, user_id,  title, description) VALUES (?, ?, ?, ?)
    `, goalID, userID, title, description)
	if err != nil {
		return 0, err
	}
	newID, _ := result.LastInsertId()
	return newID, nil
}

// UpdateEntry updates an existing learning entry.
func (service *learningStore) UpdateEntry(id int, goal_id int, user_id int, title string, description string, date time.Time, status string) error {
	_, err := service.DB.Exec(`
        UPDATE learning_entries SET title=?, description=?, date=?, status=? WHERE id=? and goal_id=? and user_id=?
    `, title, description, date, status, id, goal_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteEntry deletes a learning entry by ID.
func (service *learningStore) DeleteEntry(id int, userID int) error {
	_, err := service.DB.Exec(`
        DELETE FROM learning_entries WHERE id=? and user_id=?
    `, id, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllEntriesByGoalID returns all learning entries for a given goal ID.
func (service *learningStore) GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error) {
	var entries []models.LearningEntry
	rows, err := service.DB.Query(`
        SELECT id, goal_id, user_id, title, description, date, status FROM learning_entries WHERE goal_id=?
    `, goalID)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.LearningEntry
		err = rows.Scan(&entry.ID, &entry.GoalID, &entry.UserID, &entry.Title, &entry.Description, &entry.Date, &entry.Status)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetEntryByID returns a learning entry by ID.
func (service *learningStore) GetEntryByID(id int) (models.LearningEntry, error) {
	var entry models.LearningEntry
	err := service.DB.QueryRow(`
        SELECT id, goal_id, user_id, title, description, date, status FROM learning_entries WHERE id=?
    `, id).Scan(&entry.ID, &entry.GoalID, &entry.UserID, &entry.Title, &entry.Description, &entry.Date, &entry.Status)
	if err != nil {
		return entry, err
	}
	return entry, nil
}

// CreateFile inserts a new learning file into the database and returns its ID.
func (service *learningStore) CreateFile(entryID int, userID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error) {
	result, err := service.DB.Exec(`
        INSERT INTO learning_files (entry_id, user_id, filename, filesize, filetype, filepath) VALUES (?, ?, ?, ?, ?, ?)
    `, entryID, userID, fileName, fileSize, fileType, filePath)
	if err != nil {
		return 0, err
	}
	newID, _ := result.LastInsertId()
	return newID, nil
}

// UpdateFile updates an existing learning file.
func (service *learningStore) UpdateFile(id int, entryID int, userID int, fileName string, fileSize int64, fileType string) error {
	_, err := service.DB.Exec(`
        UPDATE learning_files SET entry_id=?, user_id=?, filename=?, filesize=?, filetype=? WHERE id=?
    `, entryID, userID, fileName, fileSize, fileType, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes a learning file by ID.
func (service *learningStore) DeleteFile(id int) error {
	_, err := service.DB.Exec(`
        DELETE FROM learning_files WHERE id=?
    `, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllFilesByGoalID returns all learning files for a given goal ID.
func (service *learningStore) GetAllFilesByEntryID(goalID int) ([]models.LearningFiles, error) {
	var files []models.LearningFiles
	rows, err := service.DB.Query(`
        SELECT id, entry_id, user_id, filename, filesize, filetype, filepath FROM learning_files WHERE entry_id=?
    `, goalID)
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		var file models.LearningFiles
		err = rows.Scan(&file.ID, &file.EntryID, &file.UserID, &file.FileName, &file.FileSize, &file.FileType, &file.FilePath)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}
	return files, nil
}

// GetFileByID returns a learning file by ID.
func (service *learningStore) GetFileByID(id int) (models.LearningFiles, error) {
	var file models.LearningFiles
	err := service.DB.QueryRow(`
        SELECT id, entry_id, user_id, filename, filesize, filetype, filepath FROM learning_files WHERE id=?
    `, id).Scan(&file.ID, &file.EntryID, &file.UserID, &file.FileName, &file.FileSize, &file.FileType, &file.FilePath)
	if err != nil {
		return file, err
	}
	return file, nil
}
