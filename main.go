package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/rs/cors"
)

type User struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	IsAccepted      bool   `json:"isAccepted"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/signup", signUpHandler)

	handler := cors.Default().Handler(mux)

	fmt.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", handler)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	log.Printf("Received POST request to /v1/signup: %+v", user)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User signed up successfully")
}
