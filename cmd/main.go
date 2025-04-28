package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"my-go-project/auth"
)

func main() {
	// Load secret key from AWS Secrets Manager
	secretsManager, err := auth.NewSecretsManager()
	if err != nil {
		log.Fatalf("Failed to initialize Secrets Manager: %v", err)
	}

	secrets, err := secretsManager.GetSecret("auth-module-secrets")
	if err != nil {
		log.Fatalf("Failed to fetch secrets: %v", err)
	}

	secretKey, exists := secrets["JWT_SECRET"]
	if !exists {
		log.Fatalf("JWT_SECRET not found in secrets")
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
