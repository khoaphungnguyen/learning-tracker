package transport

import "github.com/khoaphungnguyen/learning-tracker/internal/business"

type NetHandler struct {
	netHandler *business.LearningService
}

func NewNetHandler(netHandler *business.LearningService) *NetHandler {
	return &NetHandler{netHandler: netHandler}
}
