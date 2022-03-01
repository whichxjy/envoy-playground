package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Pong struct {
	Path     string `json:"path"`
	Cluster  string `json:"cluster"`
	Endpoint string `json:"endpoint"`
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/*", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Pong{
			Path:     c.Request().URL.Path,
			Cluster:  os.Getenv("CLUSTER"),
			Endpoint: os.Getenv("ENDPOINT"),
		})
	})

	port := os.Getenv("ENDPOINT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
