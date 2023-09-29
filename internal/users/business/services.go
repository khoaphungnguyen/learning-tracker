package userbusiness

import usermodel "github.com/khoaphungnguyen/learning-tracker/internal/users/model"

func (s *UserService) CreateUser(email string, password string, salt []byte, name string) error {
	return s.userStore.CreateUser(email, password, salt, name)
}

func (s *UserService) UpdateUser(id int, email string, name string) error {
	return s.userStore.UpdateUser(id, email, name)
}

func (s *UserService) DeleteUser(id int) error {
	return s.userStore.DeleteUser(id)
}

func (s *UserService) GetUser(id int) (usermodel.User, error) {
	return s.userStore.GetUser(id)
}

func (s *UserService) GetUserByEmail(email string) (usermodel.User, error) {
	return s.userStore.GetUserByEmail(email)
}
