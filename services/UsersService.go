package services

import (
	"bougette-backend/helper"
	"bougette-backend/models"
	"bougette-backend/repositories"
)

type UsersService struct {
	UsersRepos *repositories.UsersRepository
}

func NewUsersService(usersRepo *repositories.UsersRepository) *UsersService {
	return &UsersService{UsersRepos: usersRepo}
}

func (u *UsersService) RegisterUser(user *models.Users) error {
	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword
	return u.UsersRepos.RegisterUser(user)
}

func (u *UsersService) LoginUser(email, password string) (*models.Users, string, string, error) {
	user, err := u.UsersRepos.FindUserByEmail(email)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := helper.GenerateToken(*user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (u *UsersService) CheckUserExits(email string) (bool, error) {
	user, err := u.UsersRepos.FindUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (u *UsersService) GetUserByEmail(email string) (*models.Users, error) {
	return u.UsersRepos.FindUserByEmail(email)
}
