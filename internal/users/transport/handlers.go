package usertransport

import userbusiness "github.com/khoaphungnguyen/learning-tracker/internal/users/business"

type UserHandler struct {
	userHandler *userbusiness.UserService
	JWTKey      string
}

func NewUserHandler(userHandler *userbusiness.UserService, JWTKey string) *UserHandler {
	return &UserHandler{userHandler: userHandler, JWTKey: JWTKey}
}


