package business

import (
	"github.com/khoaphungnguyen/learning-tracker/internal/storage"
)

type LearningService struct {
	learningStore storage.LearningStore
}

func NewLearningService(learningStore storage.LearningStore) *LearningService {
	return &LearningService{learningStore: learningStore}
}
