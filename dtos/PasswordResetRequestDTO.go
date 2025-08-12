package dtos

type PasswordResetRequestDTO struct {
	Email       string `json:"email" validate:"required,email"`
	FrontendURL string `json:"frontend_url" validate:"required,url"`
}

type PasswordResetTokenDTO struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token"`
}

type PasswordResetNewPasswordDTO struct {
	Email           string `json:"email" validate:"required,email"`
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
