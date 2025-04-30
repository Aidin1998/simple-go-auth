// cmd/main.go
package main

import (
	"log"
	"os"

	"my-go-project/auth"
	"my-go-project/aws"
	"my-go-project/config"
	httpPkg "my-go-project/http"
)

func main() {
	// 1. Load configuration (PORT defaults to "80" if unset)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize AWS Secrets Manager and fetch JWT secret
	sm := aws.NewAWSSecretsManager()
	jwtSecret, err := sm.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to fetch JWT secret: %v", err)
	}

	// 3. Inject JWT secret into env for downstream compatibility
	os.Setenv("JWT_SECRET", jwtSecret)

	// 4. Set up Echo router
	router := httpPkg.SetupRouter()

	// 5. Health check / ping endpoint
	router.GET("/ping", auth.NewHandler().Ping)

	// 6. Start server on configured port
	log.Printf("Server running on port %s...\n", cfg.Port)
	log.Fatal(router.Start(":" + cfg.Port))
}
