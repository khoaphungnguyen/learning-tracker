package usermodel

import (
	"crypto/rand"
	"errors"
	"time"

	learningmodels "github.com/khoaphungnguyen/learning-tracker/internal/learning/model"
	"golang.org/x/crypto/argon2"
)

type User struct {
	ID            int                            `json:"id"`
	Email         string                         `json:"email" `
	Password      string                         `json:"password" `
	Salt          []byte                         `json:"salt"`
	Name          string                         `json:"name"`
	Role          string                         `json:"role"`
	CreatedAt     time.Time                      `json:"createdAt"`
	UpdatedAt     time.Time                      `json:"updatedAt"`
	OwnedFiles    []learningmodels.LearningFiles `json:"ownedFiles"`
	LearningGoals []learningmodels.LearningGoals `json:"learningGoals"`
}

// HashPassword takes a string as a parameter and encrypts it using argon2
func (user *User) HashPassword(password string) error {
	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	// Generate Hash
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	user.Password = string(hash)
	user.Salt = salt
	return nil
}

// CheckPassword takes a string as a parameter and compares it to the user's encrypted password
func (user *User) CheckPassword(providedPassword string) error {
	hash := argon2.IDKey([]byte(providedPassword), user.Salt, 1, 64*1024, 4, 32)
	if user.Password != string(hash) {
		return errors.New("wrong Password")
	}
	return nil
}
