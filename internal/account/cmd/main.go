package main

import (
	"net/http"
	"simple-go-auth/internal/account/config"
	"simple-go-auth/internal/account/logger"

	"github.com/brpaz/echozap"
	jwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer() *echo.Echo {
	// Load config & Zap
	config.Init()
	logger.Init() // sets logger.Logger (*zap.Logger)

	e := echo.New()

	// 1) Recovery middleware so panics become proper 5xx + a log line
	e.Use(middleware.Recover())

	// 2) Zap request‐logging middleware
	e.Use(echozap.ZapLogger(logger.Logger))

	e.Use(jwt.WithConfig(jwt.Config{
		SigningKey: []byte(config.Cfg.JWTSecret),
	}))

	// 4) Health‐check
	e.GET("/healthz", func(c echo.Context) error {
		logger.Logger.Info("Health check invoked")
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return e
}

func main() {
	e := NewServer()
	e.Logger.Fatal(e.Start(":" + config.Cfg.Port))
}
