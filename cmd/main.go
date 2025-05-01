package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"my-go-project/auth"
	"my-go-project/aws"
	"my-go-project/config"
	httpPkg "my-go-project/http"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Cognito client
	cognitoClient, err := aws.NewCognitoClient(cfg)
	if err != nil {
		log.Fatalf("Cognito init error: %v", err)
	}

	// 3. Build the auth service, handlers & middleware
	authService := auth.NewAuthServiceImpl(cognitoClient)
	authHandler := auth.NewHandler(authService)
	authMiddleware := auth.NewMiddleware(authService)

	// 4. Setup Echo router & global middleware
	router := httpPkg.SetupRouter()

	// 5. Public routes
	router.GET("/", func(c echo.Context) error {
		return c.String(200, "🚀 Welcome to BitPolaris! Please POST to /signin or /signup")
	})
	router.GET("/ping", authHandler.Ping)
	router.POST("/signup", authHandler.SignUp)
	router.POST("/signin", authHandler.SignIn)

	// 6. Protected route
	protected := router.Group("")
	protected.Use(authMiddleware)
	protected.POST("/logout", authHandler.SignOut)

	// 7. Start server
	log.Printf("Server running on port %s...\n", cfg.Port)
	log.Fatal(router.Start(":" + cfg.Port))
}
