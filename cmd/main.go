package main

import (
	"fmt"
	"log"
	"net/http"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User signed up!")
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User signed in!")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User logged out!")
}

func main() {
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
