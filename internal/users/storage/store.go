package userstorage

import (
	"database/sql"

	usermodel "github.com/khoaphungnguyen/learning-tracker/internal/users/model"
	_ "github.com/mattn/go-sqlite3" // SQLi
)

type UserStore interface {
	// Table operations
	CreateTable() error

	// User operations
	CreateUser(email string, password string, salt []byte, name string) error
	UpdateUser(id int, email string,  name string) error
	DeleteUser(id int) error
	GetUser(id int) (usermodel.User, error)
	GetUserByEmail(email string) (usermodel.User, error)
}

type userStore struct {
	DB *sql.DB
}

func NewUserStore() (*userStore, error) {
	DB, err := sql.Open("sqlite3", "./migrations/learning.db")
	if err != nil {
		panic(err)
	}
	return &userStore{DB: DB}, nil
}
