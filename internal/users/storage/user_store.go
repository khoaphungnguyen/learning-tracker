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
			name TEXT,
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
func (s *userStore) CreateUser(email string, password string, salt []byte, name string) error {
	_, err := s.DB.Exec(`
		INSERT INTO users (email, password, salt, name) VALUES (?, ?, ?, ?)
	`, email, password, salt, name)

	if err != nil {
		return err
	}
	return nil
}

// GetUser return a user's details in the database
func (s *userStore) GetUser(id int) (usermodel.User, error) {
	var user usermodel.User
	err := s.DB.QueryRow(`
		SELECT id, email, name, created_at, updated_at, role FROM users WHERE id=?
	`, id).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.Role)

	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByUserName return a user's details in the database
func (s *userStore) GetUserByEmail(email string) (usermodel.User, error) {
	var user usermodel.User
	err := s.DB.QueryRow(`
		SELECT id, password, salt FROM users WHERE email=?
	`, email).Scan(&user.ID, &user.Password, &user.Salt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates a user's details in the database
func (s *userStore) UpdateUser(id int, email string, name string) error {
	_, err := s.DB.Exec(`
		UPDATE users SET email=?, name=?, updated_at=current_timestamp  WHERE id=?
	`, email, name, id)

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID from the database
func (s *userStore) DeleteUser(id int) error {
	_, err := s.DB.Exec(`
		DELETE FROM users WHERE id=?;
	`, id)
	if err != nil {
		return err
	}
	return nil
}
