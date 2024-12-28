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

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categoryType := r.URL.Query().Get("type") // Obter o parâmetro de filtro (opcional)

	query := "SELECT id, name, type FROM categories"
	var rows *sql.Rows
	var err error

	// Filtrar por type, se fornecido
	if categoryType != "" {
		query += " WHERE type = $1"
		rows, err = config.DB.Query(query, categoryType)
	} else {
		rows, err = config.DB.Query(query)
	}

	if err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Type)
		if err != nil {
			http.Error(w, "Erro ao escanear categorias", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de categoria inválido", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = config.DB.QueryRow("SELECT id, name, type FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.Type)
	if err != nil {
		http.Error(w, "Categoria não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	if category.Type != "income" && category.Type != "expense" {
		http.Error(w, "Tipo inválido. Deve ser 'income' ou 'expense'.", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO categories (name, type) VALUES ($1, $2) RETURNING id`
	err = config.DB.QueryRow(query, category.Name, category.Type).Scan(&category.ID)
	if err != nil {
		http.Error(w, "Erro ao criar categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de categoria inválido", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	if category.Type != "income" && category.Type != "expense" {
		http.Error(w, "Tipo inválido. Deve ser 'income' ou 'expense'.", http.StatusBadRequest)
		return
	}

	query := `UPDATE categories SET name = $1, type = $2 WHERE id = $3`
	_, err = config.DB.Exec(query, category.Name, category.Type, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de categoria inválido", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM categories WHERE id = $1`
	_, err = config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Erro ao excluir categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
