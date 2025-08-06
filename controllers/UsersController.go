package controllers

import (
	"bougette-backend/common"
	"bougette-backend/dtos"
	"bougette-backend/helper"
	"bougette-backend/models"
	"bougette-backend/services"
	"bougette-backend/utilities"
	"bougette-backend/validation"
	"fmt"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	UsersService *services.UsersService
	Mailer       utilities.Mailer
}

func NewUsersController(usersService *services.UsersService, mailer utilities.Mailer) *UsersController {
	return &UsersController{UsersService: usersService, Mailer: mailer}
}

func (u *UsersController) RegisterUser(c echo.Context) error {
	request := new(dtos.UserRequestDTO)
	if err := c.Bind(request); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	exits, err := u.UsersService.CheckUserExits(request.Email)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}

	if exits {
		return common.SendBadRequestResponse(c, "User with this email already exits")
	}

	user := models.Users{
		FirstName: &request.FirstName,
		LastName:  &request.LastName,
		Gender:    &request.Gender,
		Email:     request.Email,
		Password:  request.Password,
	}

	if err := u.UsersService.RegisterUser(&user); err != nil {
		return common.SendNotFoundResponse(c, err.Error())
	}

	mailData := utilities.MailData{
		Subject: "Hello Bougette",
		Meta: struct {
			FirstName string
			LoginLink string
		}{
			FirstName: *user.FirstName,
			LoginLink: "#",
		},
	}

	if err := u.Mailer.SendViaMail(request.Email, "welcome.html", mailData); err != nil {
		fmt.Println("Email sent fail", err)
	}

	return common.SendSuccessResponse(c, "User registration successful", user)
}

func (u *UsersController) LoginUser(c echo.Context) error {
	request := new(dtos.LoginRequestDTO)
	if err := c.Bind(request); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	user, err := u.UsersService.GetUserByEmail(request.Email)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	if user == nil {
		return common.SendBadRequestResponse(c, "Invalid email or password")
	}

	if !helper.CheckPasswordHash(request.Password, user.Password) {
		return common.SendBadRequestResponse(c, "Invalid email or password")
	}

	// to generate token

	return common.SendSuccessResponse(c, "User Login successful", user)
}
