package controllers

import (
	"bougette-backend/common"
	"bougette-backend/dtos"
	"bougette-backend/services"
	"bougette-backend/validation"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	UsersService *services.UsersService
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

	return common.SendSuccessResponse(c, "User registration successful", nil)
}
