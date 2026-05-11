package storage

import (
	"encoding/json"
	"io/ioutil"

	"expense-tracker/models"
)

var (
	UsersFile        = "data/users.json"
	TransactionsFile = "data/transactions.json"
)

func ReadUsers() map[string]string {
	data, _ := ioutil.ReadFile(UsersFile)

	var users map[string]string
	json.Unmarshal(data, &users)

	if users == nil {
		return map[string]string{}
	}

	return users
}

func SaveUsers(users map[string]string) {
	data, _ := json.MarshalIndent(users, "", " ")
	ioutil.WriteFile(UsersFile, data, 0644)
}

func ReadTransactions() []models.Transaction {
	data, _ := ioutil.ReadFile(TransactionsFile)

	var tx []models.Transaction
	json.Unmarshal(data, &tx)

	if tx == nil {
		return []models.Transaction{}
	}

	return tx
}

func SaveTransactions(tx []models.Transaction) {
	data, _ := json.MarshalIndent(tx, "", " ")
	ioutil.WriteFile(TransactionsFile, data, 0644)
}