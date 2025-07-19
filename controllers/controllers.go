package controllers

import (
	"bougette-backend/dtos"
	"bougette-backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DemoController struct {
	DemoController *services.DemoServices
}

func (h *DemoController) Demo(c echo.Context) error {
	healthCheck := dtos.Demo{
		Health: true,
	}

	return c.JSON(http.StatusOK, healthCheck)
}
