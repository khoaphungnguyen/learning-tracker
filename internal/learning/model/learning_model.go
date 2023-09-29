package learningmodel

import "time"

type LearningGoals struct {
	ID        int             `json:"id"`
	UserID    int             `json:"userId"`
	Title     string          `json:"title"`
	StartDate time.Time       `json:"startDate"`
	EndDate   time.Time       `json:"endDate"`
	Entries   []LearningEntry `json:"entries"`
}

type LearningEntry struct {
	ID          int             `json:"id"`
	GoalID      int             `json:"goalId"`
	UserID      int             `json:"userId"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	Status      string          `json:"status"`
	Files       []LearningFiles `json:"files"`
}

type LearningFiles struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	EntryID   int       `json:"entryId"`
	FileName  string    `json:"fileName"`
	FileSize  int64     `json:"fileSize"`
	FileType  string    `json:"fileType"`
	FilePath  string    `json:"filePath"`
	CreatedAt time.Time `json:"createdAt"`
}
