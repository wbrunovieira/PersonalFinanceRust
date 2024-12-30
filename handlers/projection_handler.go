package handlers

import (
	"app/config"
	"app/models"
	"app/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateProjection(w http.ResponseWriter, r *http.Request) {
	var projection models.Projection

	// Lê e decodifica o JSON do corpo da requisição
	body, _ := io.ReadAll(r.Body)
	fmt.Println("Payload recebido:", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := utils.DecodeJSON(w, r, &projection); err != nil {
		return
	}

	// Verifica se os campos obrigatórios estão preenchidos
	if projection.UserID == 0 || projection.Amount == 0 || projection.CategoryID == 0 || projection.Description == "" || projection.Type == "" || projection.Date.IsZero() {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	// Validação do campo `end_month` (se for recorrente)
	if projection.IsRecurring && projection.EndMonth != nil {
		// Converte "YYYY-MM" para o tipo DATE
		endMonth, err := time.Parse("2006-01", *projection.EndMonth)
		if err != nil {
			http.Error(w, "Formato inválido para end_month. Use YYYY-MM.", http.StatusBadRequest)
			return
		}
		formattedEndMonth := endMonth.Format("2006-01-02") // Formata para DATE compatível com o PostgreSQL
		projection.EndMonth = &formattedEndMonth
	} else if projection.IsRecurring && projection.EndMonth == nil {
		http.Error(w, "Para projeções recorrentes, o campo end_month é obrigatório.", http.StatusBadRequest)
		return
	}

	// Verifica se a categoria corresponde ao tipo
	var categoryType string
	err := config.DB.QueryRow("SELECT category_type FROM categories WHERE id = $1", projection.CategoryID).Scan(&categoryType)
	if err == sql.ErrNoRows {
		http.Error(w, "Categoria não encontrada", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Erro ao verificar a categoria", http.StatusInternalServerError)
		return
	}

	if projection.Type != categoryType {
		http.Error(w, "O tipo da projeção não corresponde ao tipo da categoria", http.StatusBadRequest)
		return
	}

	// Realiza a inserção na tabela projections
	query := `
        INSERT INTO projections (user_id, amount, description, category_id, type, is_recurring, end_month, date, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id
    `
	err = config.DB.QueryRow(query, projection.UserID, projection.Amount, projection.Description, projection.CategoryID, projection.Type, projection.IsRecurring, projection.EndMonth, projection.Date).Scan(&projection.ID)
	if err != nil {
		http.Error(w, "Erro ao inserir a projeção", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(projection)
}

func GetProjections(w http.ResponseWriter, r *http.Request) {
	projectionType := r.URL.Query().Get("type")

	query := `SELECT t.id, t.user_id, t.amount, t.description, t.type, t.date, t.created_at, t.updated_at, c.name, c.category_type 
	          FROM transactions t 
	          JOIN categories c ON t.category_id = c.id`

	var rows *sql.Rows
	var err error
	if projectionType != "" {
		query += " WHERE t.type = $1"
		rows, err = config.DB.Query(query, projectionType)
	} else {
		rows, err = config.DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Erro ao buscar projeções", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projections []models.Projection
	for rows.Next() {
		var projection models.Projection
		var categoryName, categoryType string
		err = rows.Scan(&projection.ID, &projection.UserID, &projection.Amount, &projection.Description, &projection.Type, &projection.Date, &projection.CreatedAt, &projection.UpdatedAt, &categoryName, &categoryType)
		if err != nil {
			http.Error(w, "Erro ao escanear transação", http.StatusInternalServerError)
			return
		}
		projection.Category = categoryName
		projections = append(projections, projection)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projections)
}

func GetProjection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid projection ID", http.StatusBadRequest)
		return
	}

	query := `SELECT t.id, t.user_id, t.amount, t.description, t.date, t.created_at, t.updated_at, c.name 
	          FROM projections t 
	          JOIN categories c ON t.category_id = c.id 
	          WHERE t.id = $1`

	var projection models.Projection
	var categoryName string
	err = config.DB.QueryRow(query, id).Scan(&projection.ID, &projection.UserID, &projection.Amount, &projection.Description, &projection.Date, &projection.CreatedAt, &projection.UpdatedAt, &categoryName)
	if err == sql.ErrNoRows {
		http.Error(w, "Projection not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching projection", http.StatusInternalServerError)
		return
	}

	projection.Category = categoryName

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projection)
}

func UpdateProjection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de transação inválido", http.StatusBadRequest)
		return
	}

	var projection models.Projection
	err = json.NewDecoder(r.Body).Decode(&projection)
	if err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	var categoryType string
	err = config.DB.QueryRow("SELECT type FROM categories WHERE id = $1", projection.CategoryID).Scan(&categoryType)
	if err == sql.ErrNoRows {
		http.Error(w, "Categoria não encontrada", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Erro ao verificar a categoria", http.StatusInternalServerError)
		return
	}

	if projection.Type != categoryType {
		http.Error(w, "O tipo da transação não corresponde ao tipo da categoria", http.StatusBadRequest)
		return
	}

	query := `UPDATE projections 
	          SET amount = $1, description = $2, category_id = $3, type = $4, date = $5 
	          WHERE id = $6`
	_, err = config.DB.Exec(query, projection.Amount, projection.Description, projection.CategoryID, projection.Type, projection.Date, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar a transação", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteProjection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid projection ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM projections WHERE id = $1`
	_, err = config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting projection", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
