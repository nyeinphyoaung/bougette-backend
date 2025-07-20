package controllers

import (
	"bougette-backend/dtos"
	"bougette-backend/services"
	"bougette-backend/validation"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	UsersService *services.UsersService
}

func (u *UsersController) RegisterUser(c echo.Context) error {
	request := new(dtos.UserRequestDTO)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validationErrors})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "valid"})
}
