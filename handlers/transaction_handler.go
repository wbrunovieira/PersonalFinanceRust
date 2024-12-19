package handlers

import (
	"app/config"
	"app/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, "Entrada inválida. Certifique-se de que todos os campos estão corretos.", http.StatusBadRequest)
		return err
	}
	return nil
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	if err := decodeJSON(w, r, &transaction); err != nil {
		return
	}

	if transaction.UserID == 0 || transaction.Amount == 0 || transaction.CategoryID == 0 || transaction.Description == "" || transaction.Date.IsZero() {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO transactions (user_id, amount, description, category_id, date) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := config.DB.QueryRow(query, transaction.UserID, transaction.Amount, transaction.Description, transaction.CategoryID, transaction.Date).Scan(&transaction.ID)
	if err != nil {
		http.Error(w, "Erro ao inserir a transação", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	query := `SELECT t.id, t.user_id, t.amount, t.description, t.date, t.created_at, t.updated_at, c.name 
	          FROM transactions t 
	          JOIN categories c ON t.category_id = c.id`

	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		var categoryName string
		err = rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt, &categoryName)
		if err != nil {
			http.Error(w, "Error scanning transaction", http.StatusInternalServerError)
			return
		}
		transaction.Category = categoryName

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

	query := `SELECT t.id, t.user_id, t.amount, t.description, t.date, t.created_at, t.updated_at, c.name 
	          FROM transactions t 
	          JOIN categories c ON t.category_id = c.id 
	          WHERE t.id = $1`

	var transaction models.Transaction
	var categoryName string
	err = config.DB.QueryRow(query, id).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt, &categoryName)
	if err == sql.ErrNoRows {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching transaction", http.StatusInternalServerError)
		return
	}

	transaction.Category = categoryName

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

	var categoryID int
	err = config.DB.QueryRow("SELECT id FROM categories WHERE id = $1", transaction.CategoryID).Scan(&categoryID)
	if err == sql.ErrNoRows {
		http.Error(w, "Category not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error verifying category", http.StatusInternalServerError)
		return
	}

	query := `UPDATE transactions 
	          SET amount = $1, description = $2, category_id = $3, date = $4 
	          WHERE id = $5`
	_, err = config.DB.Exec(query, transaction.Amount, transaction.Description, transaction.CategoryID, transaction.Date, id)
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
