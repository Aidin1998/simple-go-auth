package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"my-go-project/auth"
	"my-go-project/config"
	"my-go-project/db"
	"my-go-project/http/health"
	"my-go-project/http/metrics"
	"my-go-project/http/ws"
)

// SetupRouter returns a fully-configured Echo instance.
func SetupRouter(h *auth.AuthHandler, authMw echo.MiddlewareFunc, logger *zap.Logger, dbConfig *config.Config) *echo.Echo {
	e := echo.New()

	// Initialize the database using InitDB
	if _, err := db.InitDB(dbConfig); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Serve HTTP/2 on TLS
	e.Server.TLSConfig.NextProtos = []string{"h2", "http/1.1"}

	// 3. Global middleware, in this order:

	// 3a. Recover from panics
	e.Use(middleware.Recover())

	// 3b. Attach X-Request-ID
	e.Use(RequestID())

	// 3c. Structured Zap logging
	e.Use(ZapLogger(logger))

	// 3d. CORS support (allow all origins by default)
	e.Use(middleware.CORS())

	// 3e. Limit body size to 2MB
	e.Use(middleware.BodyLimit("2M"))

	// 3f. Rate-limit to 10 requests/minute per IP
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// Prometheus metrics middleware
	e.Use(metrics.MetricsMiddleware())

	// 4. Debug /config endpoint
	e.GET("/config", func(c echo.Context) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, cfg)
	})

	e.GET("/ping", h.Ping)

	// Wire health endpoint
	e.GET("/health", health.HealthCheck(dbConfig, nil))

	// granular rate‐limiter
	rateLimiterStore := middleware.NewRateLimiterMemoryStore(10)
	e.POST("/signup", h.SignUp, middleware.RateLimiter(rateLimiterStore))
	e.POST("/signin", h.SignIn, middleware.RateLimiter(rateLimiterStore))
	e.POST("/refresh", h.Refresh, middleware.RateLimiter(rateLimiterStore))
	e.POST("/logout", h.SignOut, authMw)

	// Register WebSocket endpoint
	ws.RegisterWebsocket(e)

	// Expose /metrics endpoint
	metrics.MetricsHandler(e)

	return e
}
