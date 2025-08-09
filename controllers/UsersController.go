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
	"strconv"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	UsersService *services.UsersService
	Mailer       utilities.Mailer
}

func NewUsersController(usersService *services.UsersService, mailer utilities.Mailer) *UsersController {
	return &UsersController{UsersService: usersService, Mailer: mailer}
}

func (u *UsersController) GetUsers(c echo.Context) error {
	users, err := u.UsersService.GetUsers()
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "Users retrieved successfully", users)
}

func (u *UsersController) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid user ID")
	}

	user, err := u.UsersService.GetUserByID(uint(i))
	if err != nil {
		return common.SendNotFoundResponse(c, "User not found")
	}

	return common.SendSuccessResponse(c, "User retrieved successfully", user)
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

	user, accessToken, refreshToken, err := u.UsersService.LoginUser(request.Email, request.Password)
	if err != nil {
		return common.SendInternalServerErrorResponse(c, err.Error())
	}
	if user == nil {
		return common.SendBadRequestResponse(c, "Invalid email or password")
	}

	if !helper.CheckPasswordHash(request.Password, user.Password) {
		return common.SendBadRequestResponse(c, "Invalid email or password")
	}

	return common.SendSuccessResponse(c, "User Login successful", dtos.AuthResponseDTO{
		User:         dtos.MapUserToDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (u *UsersController) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid user ID")
	}

	user := new(models.Users)
	user.ID = uint(i)
	if _, err := u.UsersService.GetUserByID(user.ID); err != nil {
		return common.SendNotFoundResponse(c, "User not found")
	}

	if err := c.Bind(user); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	if err := validation.ValidateStruct(user); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	if err := u.UsersService.UpdateUser(user); err != nil {
		return common.SendNotFoundResponse(c, "User not found or update failed")
	}
	return common.SendSuccessResponse(c, "User updated successfully", user)
}

func (u *UsersController) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid user ID")
	}

	err = u.UsersService.DeleteUser(uint(i))
	if err != nil {
		return common.SendNotFoundResponse(c, "User not found")
	}

	return common.SendSuccessResponse(c, "User deleted successfully", nil)
}

func (u *UsersController) ChangePassword(c echo.Context) error {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid user ID")
	}

	request := new(dtos.ChangePasswordRequestDTO)
	if err := c.Bind(request); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	if err := u.UsersService.ChangePassword(uint(i), request); err != nil {
		return common.SendNotFoundResponse(c, "User not found or password change failed")
	}
	return common.SendSuccessResponse(c, "Password changed successfully", nil)
}
