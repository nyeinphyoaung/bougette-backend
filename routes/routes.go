package routes

import (
	"bougette-backend/controllers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitialRoute(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api/v1")

	initDemoRoutes(api, db)
}

func initDemoRoutes(e *echo.Group, db *gorm.DB) {
	controllers := controllers.UsersController{}

	// e.GET("/demo", controllers.Demo)
	e.POST("/users", controllers.RegisterUser)
}
