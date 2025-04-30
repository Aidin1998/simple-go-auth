package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"my-go-project/auth"
	"my-go-project/aws"
)

func main() {
	// Load secret key from AWS Secrets Manager
	secretsManager := aws.NewAWSSecretsManager()
	secretKey, err := secretsManager.GetJWTSecret()
	if err != nil {
		log.Fatalf("Failed to fetch JWT secret: %v", err)
	}

	// Load token expiration from environment variables (default to 1 hour)
	accessTokenExpiry := 3600 // Default 1 hour in seconds
	if val, exists := os.LookupEnv("ACCESS_TOKEN_EXPIRY"); exists {
		if parsedVal, err := strconv.Atoi(val); err == nil && parsedVal > 0 {
			accessTokenExpiry = parsedVal
		} else {
			log.Fatalf("Invalid ACCESS_TOKEN_EXPIRY value: %v", val)
		}
	}

	authService := &auth.AuthServiceImpl{
		SecretKey:         secretKey,
		AccessTokenExpiry: accessTokenExpiry,
	}
	authHandler := &auth.AuthHandler{Service: authService}
	authMiddleware := &auth.AuthMiddleware{Service: authService}

	// === New root handler for health checks and welcome ===
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("🚀 Welcome to BitPolaris! Please POST to /signin or /signup"))
	})

	// Existing endpoints
	http.HandleFunc("/signup", authHandler.SignUp)
	http.HandleFunc("/signin", authHandler.SignIn)
	http.HandleFunc("/logout", authHandler.SignOut)

	// Protected route example
	http.Handle("/protected", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected route"))
	})))

	// Use PORT environment variable or default to 80
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
