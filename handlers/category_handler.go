package handlers

import (
	"app/config"
	"app/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
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
	err = config.DB.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)
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

	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	err = config.DB.QueryRow(query, category.Name).Scan(&category.ID)
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

	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err = config.DB.Exec(query, category.Name, id)
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
