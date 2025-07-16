package main

import (
	"bougette-backend/configs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	cfg, err := configs.LoadEnv("../.env")

	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(cfg.ServerIP + ":" + cfg.ServerPort))
}
