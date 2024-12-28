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

	var categoryType string
	err := config.DB.QueryRow("SELECT type FROM categories WHERE id = $1", transaction.CategoryID).Scan(&categoryType)
	if err == sql.ErrNoRows {
		http.Error(w, "Categoria não encontrada", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Erro ao verificar a categoria", http.StatusInternalServerError)
		return
	}

	// Validar se o tipo da transação corresponde ao tipo da categoria
	if transaction.Type != categoryType {
		http.Error(w, "O tipo da transação não corresponde ao tipo da categoria", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO transactions (user_id, amount, description, category_id, type, date) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = config.DB.QueryRow(query, transaction.UserID, transaction.Amount, transaction.Description, transaction.CategoryID, transaction.Type, transaction.Date).Scan(&transaction.ID)
	if err != nil {
		http.Error(w, "Erro ao inserir a transação", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactionType := r.URL.Query().Get("type")

	query := `SELECT t.id, t.user_id, t.amount, t.description, t.type, t.date, t.created_at, t.updated_at, c.name, c.type 
	          FROM transactions t 
	          JOIN categories c ON t.category_id = c.id`

	var rows *sql.Rows
	var err error
	if transactionType != "" {
		query += " WHERE t.type = $1"
		rows, err = config.DB.Query(query, transactionType)
	} else {
		rows, err = config.DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Erro ao buscar transações", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		var categoryName, categoryType string
		err = rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Type, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt, &categoryName, &categoryType)
		if err != nil {
			http.Error(w, "Erro ao escanear transação", http.StatusInternalServerError)
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
		http.Error(w, "ID de transação inválido", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	var categoryType string
	err = config.DB.QueryRow("SELECT type FROM categories WHERE id = $1", transaction.CategoryID).Scan(&categoryType)
	if err == sql.ErrNoRows {
		http.Error(w, "Categoria não encontrada", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Erro ao verificar a categoria", http.StatusInternalServerError)
		return
	}

	if transaction.Type != categoryType {
		http.Error(w, "O tipo da transação não corresponde ao tipo da categoria", http.StatusBadRequest)
		return
	}

	query := `UPDATE transactions 
	          SET amount = $1, description = $2, category_id = $3, type = $4, date = $5 
	          WHERE id = $6`
	_, err = config.DB.Exec(query, transaction.Amount, transaction.Description, transaction.CategoryID, transaction.Type, transaction.Date, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar a transação", http.StatusInternalServerError)
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
