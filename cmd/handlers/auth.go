package handlers

import (
	"encoding/json"
	"net/http"
	"time"

"expense-tracker/models"
"expense-tracker/storage"
"expense-tracker/middleware"


	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	users := storage.ReadUsers()

	if _, exists := users[user.Username]; exists {
		http.Error(w, "User exists", 400)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	users[user.Username] = string(hash)

	storage.SaveUsers(users)

	w.WriteHeader(201)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	users := storage.ReadUsers()

	hash, exists := users[user.Username]

	if !exists ||
		bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password)) != nil {

		http.Error(w, "Invalid credentials", 401)
		return
	}

	expiration := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token, _ := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString(middleware.JwtKey)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}