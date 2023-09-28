package userstorage

import usermodel "github.com/khoaphungnguyen/learning-tracker/internal/users/model"

func (s *userStore) CreateTable() error {
	// Create the users table
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE,
			password string,
			salt BLOB,  
			full_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			role TEXT DEFAULT 'user'
		)
	`)
	if err != nil {
		return err
	}

	// Create the learning_goals table
	_, err = s.DB.Exec(`
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
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_entries (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			goal_id INTEGER,
			user_id INTEGER,
			title TEXT,
			description TEXT,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT DEFAULT 'Not Started',
			FOREIGN KEY (goal_id) REFERENCES learning_goals(id)
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create the learning_files table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS learning_files (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			entry_id INTEGER,
			user_id INTEGER,
			filename TEXT,
			filesize INTEGER,
			filetype TEXT,
			filepath TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (entry_id) REFERENCES learning_entries(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// User CRUD Methods
// CreateUser create a user's details in the database
func (s *userStore) CreateUser(email string, password string, salt []byte, fullname string) error {
	_, err := s.DB.Exec(`
		INSERT INTO users (email, password, salt, full_name) VALUES (?, ?, ?, ?)
	`, email, password, salt, fullname)

	if err != nil {
		return err
	}
	return nil
}

// GetUserByID return a user's details in the database
func (s *userStore) GetProfileByEmail(email string) (usermodel.User, error) {
	var user usermodel.User
	err := s.DB.QueryRow(`
		SELECT id, email, full_name, created_at, updated_at, role FROM users WHERE email=?
	`, email).Scan(&user.ID, &user.Email, &user.Fullname, &user.CreatedAt, &user.UpdatedAt, &user.Role)

	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByUserName return a user's details in the database
func (s *userStore) GetUserByEmail(email string) (usermodel.User, error) {
	var user usermodel.User
	err := s.DB.QueryRow(`
		SELECT email, password, salt FROM users WHERE email=?
	`, email).Scan(&user.Email, &user.Password, &user.Salt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates a user's details in the database
func (s *userStore) UpdateUser(email string, password string, fullname string) error {
	_, err := s.DB.Exec(`
		UPDATE users SET password=? full_name=?, updated_at=CURRENT_TIMESTAMP WHERE email=?
	`, password, fullname, email)

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID from the database
func (s *userStore) DeleteUser(id int, email string) error {
	_, err := s.DB.Exec(`
		DELETE FROM users WHERE id=? and email=?
	`, id, email)

	if err != nil {
		return err
	}
	return nil
}
