package health

import (
	"net/http"

	"my-go-project/aws"
	"my-go-project/config"
	"my-go-project/db"

	"github.com/labstack/echo/v4"
)

// HealthCheck returns a handler for the /health endpoint.
func HealthCheck(cfg *config.Config, sm aws.SecretsManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Initialize the database using InitDB
		dbConn, err := db.InitDB(cfg)
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"db": err.Error()})
		}

		// 1) DB ping
		sqlDB, err := dbConn.DB()
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"db": err.Error()})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"db": err.Error()})
		}

		// 2) Cognito: dummy GetUser or ListUsers call
		if _, err := sm.GetJWTSecret(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"secrets": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	}
}
