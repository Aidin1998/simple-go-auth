package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"
)

// RequestID adds a unique ID to each request and response header.
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.New().String()
			c.Response().Header().Set("X-Request-ID", id)
			c.Set("request_id", id)
			return next(c)
		}
	}
}

// ZapLogger wraps Echo’s LoggerWithConfig to include request_id.
func ZapLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("request_id", c.Get("request_id").(string)),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	})
}

// Add OpenTelemetry tracing middleware
func OpenTelemetryMiddleware(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName)
}
