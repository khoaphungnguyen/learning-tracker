package models

import "time"

type LearningEntry struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Status      string    `json:"completed"`
}

type LearningGoals struct {
	ID              int             `json:"id"`
	Title           string          `json:"title"`
	StartDate       time.Time       `json:"startdate"`
	EndDate         time.Time       `json:"enddate"`
	LearningEntries []LearningEntry `json:"LearningEntries"`
}
