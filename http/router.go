package http

import (
	"net/http"

	"my-go-project/config"

	"github.com/labstack/echo/v4"
)

// SetupRouter returns a new Echo instance.
func SetupRouter() *echo.Echo {
	e := echo.New()
	e.GET("/config", func(c echo.Context) error {
		cfg, _ := config.LoadConfig()
		return c.JSON(http.StatusOK, cfg)
	})
	return e
}
