package learningtransport

import learningbusiness "github.com/khoaphungnguyen/learning-tracker/internal/learning/business"

type LearningHandler struct {
	learningHandler *learningbusiness.LearningService
}

func NewLearningHandler(learningHandler *learningbusiness.LearningService) *LearningHandler {
	return &LearningHandler{learningHandler: learningHandler}
}
