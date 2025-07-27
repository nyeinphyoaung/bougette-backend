package routes

import (
	"bougette-backend/controllers"
	"bougette-backend/repositories"
	"bougette-backend/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitialRoute(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api/v1")

	initDemoRoutes(api, db)
}

func initDemoRoutes(e *echo.Group, db *gorm.DB) {
	usersRepos := repositories.NewUsersRepository(db)
	usersService := services.NewUsersService(usersRepos)
	usersController := controllers.NewUsersController(usersService)

	e.POST("/users", usersController.RegisterUser)
}
