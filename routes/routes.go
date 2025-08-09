package routes

import (
	"bougette-backend/controllers"
	"bougette-backend/middlewares"
	"bougette-backend/repositories"
	"bougette-backend/services"
	"bougette-backend/utilities"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitialRoute(e *echo.Echo, db *gorm.DB, mailer utilities.Mailer) {
	api := e.Group("/api/v1")

	initUsersRoutes(api, db, mailer)
}

func initUsersRoutes(e *echo.Group, db *gorm.DB, mailer utilities.Mailer) {
	usersRepos := repositories.NewUsersRepository(db)
	usersService := services.NewUsersService(usersRepos)
	usersController := controllers.NewUsersController(usersService, mailer)

	e.GET("/users", usersController.GetUsers, middlewares.IsAuthenticated)
	e.POST("/register", usersController.RegisterUser)
	e.POST("/login", usersController.LoginUser)
}
