package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
    "github.com/dgrijalva/jwt-go"
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

// CustomClaims represents custom claims for the JWT token
type CustomClaims struct {
	Username string `json:"username"`
	UserId int `json:"username"`
	jwt.StandardClaims
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/signup", signUpHandler)
	mux.HandleFunc("/v1/login", loginHandler)
	mux.HandleFunc("/v1/brackets", bracketHandler)

	handler := cors.Default().Handler(mux)

	fmt.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", handler)
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

	//if loginInfo.Username != "testuser@gmail.com" || loginInfo.Password != "testpassword" {
   	//	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    // 	return
	//}
	user := getUsersByEmail(loginInfo.Username)

	if user == nil {
		fmt.Fprint(w, "User not found")
	}

	token, err := generateJWT(user.UserId)

	//if err != nil {
    //		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
// 		return
//}
	log.Printf("Token: %+v", token)

    w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}

type Bracket struct {
	Name string `json:"name"`
}

func bracketHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request body: %s", r.Header)


	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(tokenHeader, "Bearer ") {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(tokenHeader, "Bearer ")

	log.Printf("Bearer token in request: %s", token)

	brackets := []Bracket{
		{Name: "brians-bracket"},
		{Name: "frinas-bracket"},
	}

	responseJSON, err := json.Marshal(brackets)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

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
