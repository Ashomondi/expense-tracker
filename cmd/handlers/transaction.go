package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"expense-tracker/models"
	"expense-tracker/storage"
	"expense-tracker/middleware"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := middleware.ValidateToken(r)

	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	all := storage.ReadTransactions()

	var userTx []models.Transaction

	for _, tx := range all {
		if tx.Username == username {
			userTx = append(userTx, tx)
		}
	}

	json.NewEncoder(w).Encode(userTx)
}

func AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	username, err := middleware.ValidateToken(r)

	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}

	var tx models.Transaction

	json.NewDecoder(r.Body).Decode(&tx)

	tx.Username = username
	tx.Date = time.Now().Format("2006-01-02")

	all := storage.ReadTransactions()

	all = append(all, tx)

	storage.SaveTransactions(all)

	w.WriteHeader(201)
}