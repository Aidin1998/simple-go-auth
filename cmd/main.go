// cmd/main.go
package main

import (
	"log"
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

	// 2. Fetch JWT secret via SecretsManager
	var sm aws.SecretsManager
	if cfg.Env == "dev" {
		sm = aws.NewLocalSecretsManager()
	} else {
		sm = aws.NewAWSSecretsManager(cfg.AWSRegion)
	}
	jwtSecret, err := sm.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to fetch JWT secret: %v", err)
	}
	// Store for downstream use and env-compatibility
	cfg.JWTSecret = jwtSecret
	os.Setenv("JWT_SECRET", jwtSecret)

	// 3. Setup Echo router
	router := httpPkg.SetupRouter()

	// 4. Root route (health-check & welcome)
	router.GET("/", func(c echo.Context) error {
		return c.String(200,
			"🚀 Welcome to BitPolaris! Please POST to /signin or /signup")
	})

	// 5. Ping route for smoke-tests
	router.GET("/ping", auth.NewHandler().Ping)

	// 6. Start server
	log.Printf("Server running on port %s...\n", cfg.Port)
	log.Fatal(router.Start(":" + cfg.Port))
}
