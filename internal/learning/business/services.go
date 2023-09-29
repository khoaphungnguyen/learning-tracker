package learningbusiness

import learningstorage "github.com/khoaphungnguyen/learning-tracker/internal/learning/storage"

type LearningService struct {
	learningStore learningstorage.LearningStore
}

func NewLearningService(learningStore learningstorage.LearningStore) *LearningService {
	return &LearningService{learningStore: learningStore}
}
