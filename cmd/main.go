package main

import (
	"fmt"
	"log"
	"net/http"

	"my-go-project/pkg/auth"
)

func main() {
	// Initialize AuthService and AuthHandler
	secretKey := "your-secret-key" // Replace with dynamic secret fetching
	authService := &auth.AuthServiceImpl{SecretKey: secretKey}
	authHandler := &auth.AuthHandler{Service: authService}
	authMiddleware := &auth.AuthMiddleware{Service: authService}

	http.HandleFunc("/signup", authHandler.SignUp)
	http.HandleFunc("/signin", authHandler.SignIn)
	http.HandleFunc("/logout", authHandler.SignOut)

	// Protected route example
	http.Handle("/protected", authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected route"))
	})))

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
