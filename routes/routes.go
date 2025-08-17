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
	initCategoriesRoutes(api, db)
}

func initUsersRoutes(e *echo.Group, db *gorm.DB, mailer utilities.Mailer) {
	usersRepos := repositories.NewUsersRepository(db)
	usersService := services.NewUsersService(usersRepos)
	usersController := controllers.NewUsersController(usersService, mailer)

	e.GET("/users", usersController.GetUsers, middlewares.IsAuthenticated)
	e.GET("/user/:id", usersController.GetUserByID, middlewares.IsAuthenticated)
	e.POST("/register", usersController.RegisterUser)
	e.POST("/login", usersController.LoginUser)
	e.PUT("/user/:id", usersController.UpdateUser, middlewares.IsAuthenticated)
	e.PUT("/user/:id/change-password", usersController.ChangePassword, middlewares.IsAuthenticated)
	e.DELETE("/user/:id", usersController.DeleteUser, middlewares.IsAuthenticated)
	e.POST("/forgot-password", usersController.ForgotPassword)
	e.POST("/validate-password-reset-token", usersController.ValidatePasswordResetToken)
	e.POST("/reset-password", usersController.ResetPassword)
}

func initCategoriesRoutes(e *echo.Group, db *gorm.DB) {
	categoriesRepos := repositories.NewCategoriesRepository(db)
	categoriesService := services.NewCategoriesService(categoriesRepos)
	categoriesController := controllers.NewCategoriesController(categoriesService)

	e.GET("/categories", categoriesController.GetPaginatedCategories, middlewares.IsAuthenticated)
	e.GET("/category/:id", categoriesController.GetCategoryByID, middlewares.IsAuthenticated)
	e.POST("/categories", categoriesController.CreateCategory, middlewares.IsAuthenticated)
	e.PUT("/category/:id", categoriesController.UpdateCategory, middlewares.IsAuthenticated)
	e.DELETE("/category/:id", categoriesController.DeleteCategory, middlewares.IsAuthenticated)
}
