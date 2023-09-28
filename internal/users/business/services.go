package userbusiness

import usermodel "github.com/khoaphungnguyen/learning-tracker/internal/users/model"

func (s *UserService) CreateUser(email string, password string, salt []byte, fullname string) error {
	return s.userStore.CreateUser(email, password, salt, fullname)
}

func (s *UserService) UpdateUser(email string, password string, fullname string) error{
	return s.userStore.UpdateUser(email, password, fullname)
}

func (s *UserService) DeleteUser(id int, email string) error {
	return s.userStore.DeleteUser(id, email)
}

func (s *UserService) GetProfileByEmail(email string) (usermodel.User, error) {
	return s.userStore.GetProfileByEmail(email)
}

func (s *UserService) GetUserByEmail(email string) (usermodel.User, error) {
	return s.userStore.GetUserByEmail(email)
}
