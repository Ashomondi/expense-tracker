package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Transaction struct {
	Username string  `json:"username"`
	Desc     string  `json:"desc"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Amount   float64 `json:"amount"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}