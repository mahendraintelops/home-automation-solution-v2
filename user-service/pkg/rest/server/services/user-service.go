package services

import (
	"github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/rest/server/daos"
	"github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/rest/server/models"
)

type UserService struct {
	userDao *daos.UserDao
}

func NewUserService() (*UserService, error) {
	userDao, err := daos.NewUserDao()
	if err != nil {
		return nil, err
	}
	return &UserService{
		userDao: userDao,
	}, nil
}

func (userService *UserService) CreateUser(user *models.User) (*models.User, error) {
	return userService.userDao.CreateUser(user)
}

func (userService *UserService) ListUsers() ([]*models.User, error) {
	return userService.userDao.ListUsers()
}
