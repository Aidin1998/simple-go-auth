package http

import (
	echozap "github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"my-go-project/auth"
	"my-go-project/config"
)

// SetupRouter returns a fully-configured Echo instance.
func SetupRouter(h *auth.AuthHandler, authMw echo.MiddlewareFunc) *echo.Echo {
	logger, _ := zap.NewProduction()
	e := echo.New()

	// 1. Initialize Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	// Note: call logger.Sync() on shutdown in main.go

	// 3. Global middleware, in this order:

	// 3a. Recover from panics
	e.Use(middleware.Recover())

	// 3b. Structured Zap logging
	e.Use(echozap.ZapLogger(logger))

	// 3c. CORS support (allow all origins by default)
	e.Use(middleware.CORS())

	// 3d. Limit body size to 2MB
	e.Use(middleware.BodyLimit("2M"))

	// 3e. Rate-limit to 10 requests/minute per IP
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// 4. Debug /config endpoint
	e.GET("/config", func(c echo.Context) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, cfg)
	})

	e.GET("/ping", h.Ping)

	// granular rate‐limiter
	rateLimiterStore := middleware.NewRateLimiterMemoryStore(10)
	e.POST("/signup", h.SignUp, middleware.RateLimiter(rateLimiterStore))
	e.POST("/signin", h.SignIn, middleware.RateLimiter(rateLimiterStore))
	e.POST("/refresh", h.Refresh, middleware.RateLimiter(rateLimiterStore))
	e.POST("/logout", h.SignOut, authMw)
	return e
}
