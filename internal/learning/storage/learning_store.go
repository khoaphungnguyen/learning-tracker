package learningstorage

import (
	"time"

	learningmodel "github.com/khoaphungnguyen/learning-tracker/internal/learning/model"
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
		UPDATE learning_goals SET title=?, startdate=?, enddate=? WHERE user_id=? and id=?
	`, title, startDate, endDate, userID, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGoal deletes a learning goal by ID.
func (service *learningStore) DeleteGoal(id int, userID int) error {
	_, err := service.DB.Exec(`
		DELETE FROM learning_goals WHERE id=? and user_id=?
	`, id, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllGoalsByUserID returns all learning goals for a given user ID.
func (service *learningStore) GetAllGoalsByUserID(userID int) ([]learningmodel.LearningGoals, error) {
	var goals []learningmodel.LearningGoals
	rows, err := service.DB.Query(`
		SELECT id, title, startdate, enddate FROM learning_goals WHERE user_id=?
	`, userID)
	if err != nil {
		return goals, err
	}
	defer rows.Close()

	for rows.Next() {
		var goal learningmodel.LearningGoals
		err = rows.Scan(&goal.ID, &goal.Title, &goal.StartDate, &goal.EndDate)
		if err != nil {
			return goals, err
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

// GetGoalByID returns a learning goal by ID.
func (service *learningStore) GetGoalByID(id int) (learningmodel.LearningGoals, error) {
	var goal learningmodel.LearningGoals
	err := service.DB.QueryRow(`
		SELECT id, title, startdate, enddate FROM learning_goals WHERE id=?
	`, id).Scan(&goal.ID, &goal.Title, &goal.StartDate, &goal.EndDate)
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
func (service *learningStore) UpdateEntry(id int, user_id int, title string, description string, status string) error {
	_, err := service.DB.Exec(`
        UPDATE learning_entries SET title=?, description=?, date=current_timestamp, status=? WHERE id=? and user_id=?
    `, title, description, status, id, user_id)
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
func (service *learningStore) GetAllEntriesByGoalID(goalID int) ([]learningmodel.LearningEntry, error) {
	var entries []learningmodel.LearningEntry
	rows, err := service.DB.Query(`
        SELECT id, title, description, date, status FROM learning_entries WHERE goal_id=?
    `, goalID)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry learningmodel.LearningEntry
		err = rows.Scan(&entry.ID, &entry.Title, &entry.Description, &entry.Date, &entry.Status)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetEntryByID returns a learning entry by ID.
func (service *learningStore) GetEntryByID(id int) (learningmodel.LearningEntry, error) {
	var entry learningmodel.LearningEntry
	err := service.DB.QueryRow(`
        SELECT id, title, description, date, status FROM learning_entries WHERE id=?
    `, id).Scan(&entry.ID, &entry.Title, &entry.Description, &entry.Date, &entry.Status)
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
func (service *learningStore) UpdateFile(id int, userID int, fileName string, fileSize int64, fileType string, filePath string) error {
	_, err := service.DB.Exec(`
        UPDATE learning_files SET filename=?, filesize=?, filetype=?, filePath=? WHERE id=? and user_id=?
    `, fileName, fileSize, fileType, filePath, id, userID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes a learning file by ID.
func (service *learningStore) DeleteFile(id int, userID int) error {
	_, err := service.DB.Exec(`
        DELETE FROM learning_files WHERE id=? and user_id=?
    `, id, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetAllFilesByGoalID returns all learning files for a given goal ID.
func (service *learningStore) GetAllFilesByEntryID(entryID int) ([]learningmodel.LearningFiles, error) {
	var files []learningmodel.LearningFiles
	rows, err := service.DB.Query(`
        SELECT id, filename, filesize, filetype, filepath FROM learning_files WHERE entry_id=?
    `, entryID)
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		var file learningmodel.LearningFiles
		err = rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.FileType, &file.FilePath)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}
	return files, nil
}

// GetFileByID returns a learning file by ID.
func (service *learningStore) GetFileByID(id int) (learningmodel.LearningFiles, error) {
	var file learningmodel.LearningFiles
	err := service.DB.QueryRow(`
        SELECT id, filename, filesize, filetype, filepath FROM learning_files WHERE id=?
    `, id).Scan(&file.ID, &file.FileName, &file.FileSize, &file.FileType, &file.FilePath)
	if err != nil {
		return file, err
	}
	return file, nil
}
