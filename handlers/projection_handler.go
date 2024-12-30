package handlers

import (
	"app/config"
	"app/models"
	"app/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateProjection(w http.ResponseWriter, r *http.Request) {
	var projection models.Projection

	if err := utils.DecodeJSON(w, r, &projection); err != nil {
		return
	}

	if projection.UserID == 0 || projection.Amount == 0 || projection.CategoryID == 0 || projection.Description == "" || projection.Date.IsZero() {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "O tipo da transação não corresponde ao tipo da categoria", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO projections (user_id, amount, description, category_id, type, date) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = config.DB.QueryRow(query, projection.UserID, projection.Amount, projection.Description, projection.CategoryID, projection.Type, projection.Date).Scan(&projection.ID)
	if err != nil {
		http.Error(w, "Erro ao inserir a transação", http.StatusInternalServerError)
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
