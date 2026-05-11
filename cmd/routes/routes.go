package routes

import (
	"net/http"

	"expense-tracker/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(frontendDir string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	r.HandleFunc("/transactions", handlers.GetTransactionsHandler).Methods("GET")
	r.HandleFunc("/transactions", handlers.AddTransactionHandler).Methods("POST")

	fs := http.FileServer(http.Dir(frontendDir))
	r.PathPrefix("/").Handler(fs)

	return r
}