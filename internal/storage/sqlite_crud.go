package storage

import (
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func (service *learningStore) CreateTable() error {
	// Create the users table
	_, err := service.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			password_hash TEXT,
			salt BLOB,  
			first_name TEXT,
			last_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			role TEXT DEFAULT 'user'
		)
	`)
	if err != nil {
		return err
	}

	// Create the learning_goals table
	_, err = service.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_goals (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			startdate DATETIME,
			enddate DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create the learning_entries table
	_, err = service.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_entries (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			goal_id INTEGER,
			title TEXT,
			description TEXT,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT DEFAULT 'Not Started',
			FOREIGN KEY (goal_id) REFERENCES learning_goals(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create the learning_files table
	_, err = service.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_files (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			goal_id INTEGER,
			filename TEXT,
			filesize INTEGER,
			filetype TEXT,
			filepath TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			owner_id INTEGER,
			FOREIGN KEY (goal_id) REFERENCES learning_goals(id),
			FOREIGN KEY (owner_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// User CRUD Methods
// CreateUser create a user's details in the database
func (service *learningStore) CreateUser(username string, passwordHash string, hash []byte, firstName string, lastName string) (int64, error) {
	result, err := service.DB.Exec(`
		INSERT INTO users (username, password_hash, hash, first_name, last_name, role) VALUES (?, ?, ?, ?, ?)
	`, username, passwordHash, hash, firstName, lastName)

	if err != nil {
		return 0, err
	}

	newID, err := result.LastInsertId()
	return newID, err
}

// GetUserByID return a user's details in the database
func (service *learningStore) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := service.DB.QueryRow(`
		SELECT id, username, first_name, last_name, created_at, updated_at, role FROM users WHERE id=?
	`, id).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.Role)

	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByUserName return a user's details in the database
func (service *learningStore) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := service.DB.QueryRow(`
		SELECT id, username, first_name, last_name, created_at, updated_at, role FROM users WHERE username=?
	`, username).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.Role)

	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates a user's details in the database
func (service *learningStore) UpdateUser(id int, username string, firstName string, lastName string) error {
	_, err := service.DB.Exec(`
		UPDATE users SET username=?, first_name=?, last_name=?, updated_at=CURRENT_TIMESTAMP WHERE id=?
	`, username, firstName, lastName, id)

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID from the database
func (service *learningStore) DeleteUser(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM users WHERE id=?
	`, id)

	if err != nil {
		return err
	}
	return nil
}

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
func (service *learningStore) CreateEntry(goalID int, title string, description string) (int64, error) {
	result, err := service.DB.Exec(`
        INSERT INTO learning_entries (goal_id, title, description) VALUES (?, ?, ?)
    `, goalID, title, description)
	if err != nil {
		return 0, err
	}
	newID, _ := result.LastInsertId()
	return newID, nil
}

// UpdateEntry updates an existing learning entry.
func (service *learningStore) UpdateEntry(id int, title string, description string, date time.Time, status string) error {
	_, err := service.DB.Exec(`
        UPDATE learning_entries SET title=?, description=?, date=?, status=? WHERE id=?
    `, title, description, date, status, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteEntry deletes a learning entry by ID.
func (service *learningStore) DeleteEntry(id int) error {
	_, err := service.DB.Exec(`
        DELETE FROM learning_entries WHERE id=?
    `, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllEntriesByGoalID returns all learning entries for a given goal ID.
func (service *learningStore) GetAllEntriesByGoalID(goalID int) ([]models.LearningEntry, error) {
	var entries []models.LearningEntry
	rows, err := service.DB.Query(`
        SELECT id, goal_id, title, description, date, status FROM learning_entries WHERE goal_id=?
    `, goalID)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.LearningEntry
		err = rows.Scan(&entry.ID, &entry.LearningGoalID, &entry.Title, &entry.Description, &entry.Date, &entry.Status)
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
        SELECT id, goal_id, title, description, date, status FROM learning_entries WHERE id=?
    `, id).Scan(&entry.ID, &entry.LearningGoalID, &entry.Title, &entry.Description, &entry.Date, &entry.Status)
	if err != nil {
		return entry, err
	}
	return entry, nil
}

// CreateFile inserts a new learning file into the database and returns its ID.
func (service *learningStore) CreateFile(goalID int, ownerID int, fileName string, fileSize int64, fileType string, filePath string) (int64, error) {
	result, err := service.DB.Exec(`
        INSERT INTO learning_files (goal_id, owner_id, filename, filesize, filetype, filepath) VALUES (?, ?, ?, ?, ?, ?)
    `, goalID, ownerID, fileName, fileSize, fileType, filePath)
	if err != nil {
		return 0, err
	}
	newID, _ := result.LastInsertId()
	return newID, nil
}

// UpdateFile updates an existing learning file.
func (service *learningStore) UpdateFile(id int, goalID int, ownerID int, fileName string, fileSize int64, fileType string) error {
	_, err := service.DB.Exec(`
        UPDATE learning_files SET goal_id=?, owner_id=?, filename=?, filesize=?, filetype=? WHERE id=?
    `, goalID, ownerID, fileName, fileSize, fileType, id)
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
func (service *learningStore) GetAllFilesByGoalID(goalID int) ([]models.LearningFiles, error) {
	var files []models.LearningFiles
	rows, err := service.DB.Query(`
        SELECT id, goal_id, owner_id, filename, filesize, filetype, filepath FROM learning_files WHERE goal_id=?
    `, goalID)
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		var file models.LearningFiles
		err = rows.Scan(&file.ID, &file.LearningGoalID, &file.OwnerID, &file.FileName, &file.FileSize, &file.FileType, &file.FilePath)
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
        SELECT id, goal_id, owner_id, filename, filesize, filetype, filepath FROM learning_files WHERE id=?
    `, id).Scan(&file.ID, &file.LearningGoalID, &file.OwnerID, &file.FileName, &file.FileSize, &file.FileType, &file.FilePath)
	if err != nil {
		return file, err
	}
	return file, nil
}
