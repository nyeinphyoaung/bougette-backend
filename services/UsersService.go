package services

import (
	"bougette-backend/dtos"
	"bougette-backend/helper"
	"bougette-backend/models"
	"bougette-backend/repositories"
	"errors"
	"fmt"
	"time"
)

type UsersService struct {
	UsersRepos *repositories.UsersRepository
}

func NewUsersService(usersRepo *repositories.UsersRepository) *UsersService {
	return &UsersService{UsersRepos: usersRepo}
}

func (u *UsersService) GetUsers() ([]models.Users, error) {
	return u.UsersRepos.GetUsers()
}

func (u *UsersService) GetUserByID(id uint) (*models.Users, error) {
	return u.UsersRepos.GetUserByID(id)
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

	if user == nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := helper.GenerateToken(*user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (u *UsersService) UpdateUser(user *models.Users) error {
	if user.Password != "" {
		hashPassword, err := helper.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashPassword
	}

	return u.UsersRepos.UpdateUser(user)
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

func (u *UsersService) DeleteUser(id uint) error {
	return u.UsersRepos.DeleteUser(id)
}

func (u *UsersService) ChangePassword(id uint, request *dtos.ChangePasswordRequestDTO) error {
	user, err := u.UsersRepos.GetUserByID(id)
	if err != nil {
		return err
	}

	if !helper.CheckPasswordHash(request.CurrentPassword, user.Password) {
		return fmt.Errorf("current password is incorrect")
	}

	hashPassword, err := helper.HashPassword(request.NewPassword)
	if err != nil {
		return nil
	}

	user.Password = hashPassword
	return u.UsersRepos.UpdateUser(user)
}

func (u *UsersService) ForgotPassword(email string) (string, error) {
	user, err := u.UsersRepos.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found for password reset")
	}

	token := helper.GeneratePasswordToken()
	if token == "" {
		return "", errors.New("failed to generate token")
	}

	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 15), // 15 min
		Used:      false,
	}

	if err := u.UsersRepos.InvalidateOldPasswordResets(user.ID); err != nil {
		fmt.Printf("Warning: Could not invalidate old password resets for user %d: %v\n", user.ID, err)
	}

	if err := u.UsersRepos.CreatePasswordReset(passwordReset); err != nil {
		return "", err
	}

	return token, nil
}

func (u *UsersService) ValidatePasswordResetToken(token dtos.PasswordResetTokenDTO) error {
	passwordReset, err := u.UsersRepos.GetPasswordResetByTokenAndEmail(token.Token, token.Email)
	if err != nil {
		return errors.New("invalid or expired password reset token")
	}

	if passwordReset.ExpiredAt.Before(time.Now()) {
		return errors.New("password reset token has expired")
	}
	if passwordReset.Used {
		return errors.New("password reset token has already been used")
	}
	return nil
}

func (u *UsersService) ResetPassword(request *dtos.PasswordResetNewPasswordDTO) error {
	user, err := u.UsersRepos.FindUserByEmail(request.Email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	passwordReset, err := u.UsersRepos.GetPasswordResetByTokenAndEmail(request.Token, request.Email)
	if err != nil {
		return errors.New("invalid or expired password reset token")
	}

	if passwordReset.UserID != user.ID {
		return errors.New("invalid password reset token for this user")
	}

	if passwordReset.ExpiredAt.Before(time.Now()) {
		return errors.New("password reset token has expired")
	}
	if passwordReset.Used {
		return errors.New("password reset token has already been used")
	}

	hashedPassword, err := helper.HashPassword(request.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	err = u.UsersRepos.ResetUserPasswordAndMarkTokenUsed(user, hashedPassword, passwordReset)
	if err != nil {
		return errors.New("failed to reset password, please try again")
	}

	return nil
}
