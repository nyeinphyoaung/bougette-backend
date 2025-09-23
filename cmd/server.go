package main

import (
	"bougette-backend/configs"
	"bougette-backend/middlewares"
	"bougette-backend/repositories"
	"bougette-backend/routes"
	"bougette-backend/services"
	"bougette-backend/utilities"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	s := middlewares.NewStats()
	e.Use(s.Process)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/stats", s.Handle)

	cfg := configs.Envs

	if err := cfg.ConnectDB(); err != nil {
		e.Logger.Fatal(err)
	}

	cfg.InitializedDB()
	mailer := utilities.NewMailer()

	e.Use(middlewares.ServerHeader, middleware.Logger())
	routes.InitialRoute(e, cfg.DB, mailer, services.NewNotificationsService(repositories.NewNotificationsRepos(cfg.DB)))
	e.Logger.Fatal(e.Start(cfg.ServerIP + ":" + cfg.ServerPort))
}
