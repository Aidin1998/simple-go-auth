// cmd/main.go
package main

import (
	"log"
	"net/http"

	"os"

	"github.com/labstack/echo/v4"

	"my-go-project/auth"
	"my-go-project/aws"
	"my-go-project/config"
	httpPkg "my-go-project/http"
)

func main() {
	// 1. Load config (PORT defaults to "80")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Get JWT secret
	sm := aws.NewAWSSecretsManager()
	jwtSecret, err := sm.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to fetch JWT secret: %v", err)
	}
	os.Setenv("JWT_SECRET", jwtSecret)

	// 3. Setup Echo
	router := httpPkg.SetupRouter()

	// 4. Root route for health-checks & welcome
	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK,
			"🚀 Welcome to BitPolaris! Please POST to /signin or /signup")
	})

	// 5. Ping route for smoke-tests
	router.GET("/ping", auth.NewHandler().Ping)

	// 6. Start server
	log.Printf("Server running on port %s...\n", cfg.Port)
	log.Fatal(router.Start(":" + cfg.Port))
}
