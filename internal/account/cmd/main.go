package main

import (
	"simple-go-auth/internal/account/config"
	"simple-go-auth/internal/account/logger"
	"simple-go-auth/internal/users/auth"

	"go.uber.org/zap"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ZapEchoLogger struct {
	*zap.SugaredLogger
}

func NewZapEchoLogger(sugar *zap.SugaredLogger) *ZapEchoLogger {
	return &ZapEchoLogger{SugaredLogger: sugar}
}

// Implement echo.Logger interface methods for ZapEchoLogger
func (l *ZapEchoLogger) Debugj(j map[string]interface{}) {
	l.Debugw("debug", "data", j)
}

func (l *ZapEchoLogger) Infoj(j map[string]interface{}) {
	l.Infow("info", "data", j)
}

func (l *ZapEchoLogger) Warnj(j map[string]interface{}) {
	l.Warnw("warn", "data", j)
}

// Removed misplaced line

// NewServer constructs and returns an Echo instance with all middleware & routes.
func NewServer() *echo.Echo {
	config.Init()
	logger.Init()

	e := echo.New()
	e.Logger = NewZapEchoLogger(logger.Logger.Sugar())

	// JWT middleware stub; uses simple-go-auth to get secret
	e.Use(middleware.JWT(middleware.JWTConfig{
		SigningKey: []byte(auth.GetJWTSecret()),
	}))

	// Health check
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
