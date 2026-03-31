package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	usersFile        = "users.json"
	transactionsFile = "transactions.json"
	jwtKey           = []byte("supersecretkey")
)

type (
	User        struct{ Username, Password string }
	Transaction struct {
		Username, Desc, Category, Date string
		Amount                         float64
	}
)
type Claims struct {
	Username string
	jwt.RegisteredClaims
}

func readUsers() map[string]string {
	data, _ := ioutil.ReadFile(usersFile)
	var m map[string]string
	json.Unmarshal(data, &m)
	if m == nil {
		return map[string]string{}
	}
	return m
}

func saveUsers(u map[string]string) {
	b, _ := json.MarshalIndent(u, "", " ")
	ioutil.WriteFile(usersFile, b, 0644)
}

func readTransactions() []Transaction {
	data, _ := ioutil.ReadFile(transactionsFile)
	var t []Transaction
	json.Unmarshal(data, &t)
	if t == nil {
		return []Transaction{}
	}
	return t
}

func saveTransactions(t []Transaction) {
	b, _ := json.MarshalIndent(t, "", " ")
	ioutil.WriteFile(transactionsFile, b, 0644)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	users := readUsers()
	if _, exists := users[u.Username]; exists {
		http.Error(w, "User exists", 400)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	users[u.Username] = string(hash)
	saveUsers(users)
	w.WriteHeader(201)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	users := readUsers()
	hash, exists := users[u.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password)) != nil {
		http.Error(w, "Invalid", 401)
		return
	}
	exp := time.Now().Add(24 * time.Hour)
	claims := &Claims{Username: u.Username, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(exp)}}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func validateToken(r *http.Request) (string, error) {
	t := r.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(t, &Claims{}, func(t *jwt.Token) (interface{}, error) { return jwtKey, nil })
	if err != nil {
		return "", err
	}
	claims := token.Claims.(*Claims)
	return claims.Username, nil
}

func getTxHandler(w http.ResponseWriter, r *http.Request) {
	u, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}
	all := readTransactions()
	userTx := []Transaction{}
	for _, t := range all {
		if t.Username == u {
			userTx = append(userTx, t)
		}
	}
	json.NewEncoder(w).Encode(userTx)
}

func addTxHandler(w http.ResponseWriter, r *http.Request) {
	u, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}
	var t Transaction
	json.NewDecoder(r.Body).Decode(&t)
	t.Username = u
	t.Date = time.Now().Format("2006-01-02")
	all := readTransactions()
	all = append(all, t)
	saveTransactions(all)
	w.WriteHeader(201)
}

func main() {
	if _, err := os.Stat(usersFile); os.IsNotExist(err) {
		ioutil.WriteFile(usersFile, []byte("{}"), 0644)
	}
	if _, err := os.Stat(transactionsFile); os.IsNotExist(err) {
		ioutil.WriteFile(transactionsFile, []byte("[]"), 0644)
	}
	r := mux.NewRouter()
	r.HandleFunc("/signup", signupHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/transactions", getTxHandler).Methods("GET")
	r.HandleFunc("/transactions", addTxHandler).Methods("POST")
	fmt.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
