package userbusiness

import userstorage "github.com/khoaphungnguyen/learning-tracker/internal/users/storage"

type UserService struct {
	userStore userstorage.UserStore
}

func NewUserService(userStore userstorage.UserStore) *UserService {
	return &UserService{userStore: userStore}
}
