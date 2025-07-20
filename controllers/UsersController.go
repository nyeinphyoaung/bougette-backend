package controllers

import (
	"bougette-backend/dtos"
	"bougette-backend/services"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	UsersService *services.UsersService
}

func (u *UsersController) RegisterUser(c echo.Context) error {
	request := new(dtos.UserRequestDTO)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	fmt.Println(request)
	return c.String(http.StatusOK, "good request")
}
