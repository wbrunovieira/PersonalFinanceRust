package handlers

import (
	"app/config"
	"app/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var userID int
	err = config.DB.QueryRow("SELECT id FROM users WHERE id = $1", transaction.UserID).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		return
	}

	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	query := `INSERT INTO transactions (user_id, amount, description, category, date) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = config.DB.QueryRow(query, transaction.UserID, transaction.Amount, transaction.Description, transaction.Category, transaction.Date).Scan(&transaction.ID)
	if err != nil {
		http.Error(w, "Error inserting transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, user_id, amount, description, category, date, created_at, updated_at FROM transactions")
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Category, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			http.Error(w, "Error scanning transaction", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	query := `SELECT id, user_id, amount, description, category, date, created_at, updated_at FROM transactions WHERE id = $1`
	err = config.DB.QueryRow(query, id).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Category, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE transactions SET amount = $1, description = $2, category = $3, date = $4 WHERE id = $5`
	_, err = config.DB.Exec(query, transaction.Amount, transaction.Description, transaction.Category, transaction.Date, id)
	if err != nil {
		http.Error(w, "Error updating transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM transactions WHERE id = $1`
	_, err = config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
