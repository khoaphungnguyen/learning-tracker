package transport

import "github.com/khoaphungnguyen/learning-tracker/internal/business"

type NetHandler struct {
	netHandler *business.LearningService
	JWTKey string
}

func NewNetHandler(netHandler *business.LearningService, JWTKey string) *NetHandler {
	return &NetHandler{netHandler: netHandler, JWTKey: JWTKey}
}
