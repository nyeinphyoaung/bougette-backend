package dtos

import "bougette-backend/models"

type UserRequestDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Gender    string `json:"gender" validate:"oneof=male female prefer_not_to"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponseDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
}

type AuthResponseDTO struct {
	User         UserResponseDTO `json:"user"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
}

func MapUserToDTO(user *models.Users) UserResponseDTO {
	return UserResponseDTO{
		ID:        user.ID,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     user.Email,
		Gender:    *user.Gender,
	}
}
