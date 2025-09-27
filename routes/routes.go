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

func InitialRoute(e *echo.Echo, db *gorm.DB, mailer utilities.Mailer, notificationsService *services.NotificationsService) {
	api := e.Group("/api/v1")

	initUsersRoutes(api, db, mailer, notificationsService)
	initCategoriesRoutes(api, db)
	initBudgetsRoutes(api, db)
	initNotificationsRoutes(api, db)
	initWalletRoutes(api, db)
	initWebsocketRoutes(api)
}

func initUsersRoutes(e *echo.Group, db *gorm.DB, mailer utilities.Mailer, notificationsService *services.NotificationsService) {
	usersRepos := repositories.NewUsersRepository(db)
	usersService := services.NewUsersService(usersRepos)
	usersController := controllers.NewUsersController(usersService, mailer, notificationsService)

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

func initBudgetsRoutes(e *echo.Group, db *gorm.DB) {
	budgetsRepos := repositories.NewBudgetsRepository(db)
	budgetsService := services.NewBudgetsService(budgetsRepos)
	categoriesRepos := repositories.NewCategoriesRepository(db)
	categoriesService := services.NewCategoriesService(categoriesRepos)
	budgetsController := controllers.NewBudgetsController(budgetsService, categoriesService)

	e.POST("/budgets", budgetsController.CreateBudgets, middlewares.IsAuthenticated)
	e.GET("/budgets", budgetsController.GetPaginatedBudgets, middlewares.IsAuthenticated)
	e.PATCH("/budgets/:id", budgetsController.UpdateBudget, middlewares.IsAuthenticated)
	e.DELETE("/budgets/:id", budgetsController.DeleteBudget, middlewares.IsAuthenticated)
}

func initNotificationsRoutes(e *echo.Group, db *gorm.DB) {
	notificationsRepos := repositories.NewNotificationsRepos(db)
	notificationsService := services.NewNotificationsService(notificationsRepos)
	notificationsController := controllers.NewNotificationsController(notificationsService)

	e.GET("/notifications/:user_id", notificationsController.GetNotificationsByUserID, middlewares.IsAuthenticated)
	e.PUT("/notifications/:notification_id/read", notificationsController.MarkNotificationAsRead, middlewares.IsAuthenticated)
	e.PUT("/notifications/:user_id/mark-all-read", notificationsController.MarkAllNotificationsAsRead, middlewares.IsAuthenticated)
	e.DELETE("/notifications/:notification_id", notificationsController.DeleteNotification, middlewares.IsAuthenticated)
	e.DELETE("/notifications/:user_id/clear", notificationsController.ClearAllNotifications, middlewares.IsAuthenticated)
}

func initWalletRoutes(e *echo.Group, db *gorm.DB) {
	walletRepos := repositories.NewWalletRepository(db)
	walletService := services.NewWalletService(walletRepos)
	walletController := controllers.NewWalletController(walletService)

	e.POST("/wallets", walletController.CreateWallet, middlewares.IsAuthenticated)
	e.GET("/wallets/default", walletController.GenerateDefaultWallet, middlewares.IsAuthenticated)
}

func initWebsocketRoutes(e *echo.Group) {
	e.GET("/ws/:user_id", controllers.HandleWebSocket)
}
