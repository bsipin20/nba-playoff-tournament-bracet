package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
    "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type UserSignUp struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	IsAccepted      bool   `json:"isAccepted"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var secretKey = []byte("your_secret_key")

// CustomClaims represents custom claims for the JWT token
type CustomClaims struct {
	Username string `json:"username"`
	UserId int `json:"username"`
	jwt.StandardClaims
}

func main() {
	router := mux.NewRouter()

    router.HandleFunc("/v1/signup", signUpHandler).Methods("POST")
    router.HandleFunc("/v1/login", loginHandler).Methods("POST")
	router.HandleFunc("/v1/brackets/{id}", bracketHandler).Methods("GET")


    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowCredentials: true,
    })

    handler := c.Handler(router)
    log.Fatal(http.ListenAndServe(":8080", handler))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user UserSignUp
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	log.Printf("Received POST request to /v1/signup: %+v", user)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User signed up successfully")
}

type UserDBRecord struct {
	UserId string
	Name string
	Password string
	Email string
}

var GlobalUsers = []UserDBRecord{
	{
		UserId: "1",
		Name: "Brian Sip",
		Password: "password123",
		Email: "bsipin@gmail.com",
	},
	{
		UserId: "2",
		Name: "Frina",
		Password: "password",
		Email: "frina.lin@gmail.com",
	},
}

func getUserById(userId string) *UserDBRecord {
	for _, user := range GlobalUsers {
		if user.UserId == userId {
			return &user
		}
	}
	return nil
}

func getUsersByEmail(email string) *UserDBRecord {
	for _, user := range GlobalUsers {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Sample user credentials (you might fetch this from a database)
	validUser := User{
		ID:       "1",
		Username: "user",
		Password: "password",
	}

	// Decode request body to get user credentials
	var userCreds User
	err := json.NewDecoder(r.Body).Decode(&userCreds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if user credentials are valid
	if userCreds.Username != validUser.Username || userCreds.Password != validUser.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": validUser.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send the token in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginInfo UserLogin
	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	log.Printf("Request %+v", loginInfo.Username)

	user := getUsersByEmail(loginInfo.Username)

	if user == nil {
		fmt.Fprint(w, "User not found")
	}

	token, err := generateJWT(user.UserId)

	log.Printf("Token: %+v", token)

    w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"token": "%s", "userId": "%s"}`, token, user.UserId)
}

type Bracket struct {
	Name string `json:"name"`
}

// CustomClaims struct to parse JWT claims
type ValidCustomClaims struct {
    UserId string `json:"user_id"`
    jwt.StandardClaims
}

func extractUserIDFromPath(path string) string {
    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        return ""
    }
    return parts[len(parts)-1]
}


func bracketHandler(w http.ResponseWriter, r *http.Request) {

	brackets := []Bracket{
		{Name: "brians-bracket"},
		{Name: "frinas-bracket"},
	}

	responseJSON, err := json.Marshal(brackets)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("Response: %s", responseJSON)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func generateJWT(username string) (string, error) {
	claims := CustomClaims{
		Username: username,
		UserId: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		log.Println("Error generating JWT token:", err)
		return "", err
	}

	return tokenString, nil
}
