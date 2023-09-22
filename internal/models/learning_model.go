package models

import "time"

type User struct {
	ID            int             `json:"id"`
	Username      string          `json:"username"`
	PasswordHash  string          `json:"-"`
	Salt          []byte          `json:"salt"`
	FirstName     string          `json:"firstName"`
	LastName      string          `json:"lastName"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	Role          string          `json:"role"`
	OwnedFiles    []LearningFiles `json:"ownedFiles"`
	LearningGoals []LearningGoals `json:"learningGoals"` // New field
}

type LearningEntry struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Date           time.Time `json:"date"`
	Status         string    `json:"completed"`
	LearningGoalID int       `json:"learningGoalId"`
}

type LearningGoals struct {
	ID              int             `json:"id"`
	UserID          int             `json:"userId"` // New field
	Title           string          `json:"title"`
	StartDate       time.Time       `json:"startDate"`
	EndDate         time.Time       `json:"endDate"`
	LearningEntries []LearningEntry `json:"learningEntries"`
	LearningFiles   []LearningFiles `json:"learningFiles"`
}

type LearningFiles struct {
	ID             int       `json:"id"`
	FileName       string    `json:"fileName"`
	FileSize       int64     `json:"fileSize"`
	FileType       string    `json:"fileType"`
	FilePath       string    `json:"filePath"`
	LearningGoalID int       `json:"learningGoalId"`
	CreatedAt      time.Time `json:"createdAt"`
	OwnerID        int       `json:"ownerId"`
}
