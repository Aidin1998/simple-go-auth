package http

import "github.com/labstack/echo/v4"

// SetupRouter returns a new Echo instance.
func SetupRouter() *echo.Echo {
	e := echo.New()
	return e
}
