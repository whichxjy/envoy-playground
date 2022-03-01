package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/*", func(c echo.Context) error {
		if token := c.Request().Header.Get("Token"); token != "hi" {
			return c.String(http.StatusForbidden, "Forbidden")
		}

		return c.JSON(http.StatusOK, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
