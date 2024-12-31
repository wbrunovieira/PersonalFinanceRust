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

	body, _ := io.ReadAll(r.Body)
	fmt.Println("Payload recebido:", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := utils.DecodeJSON(w, r, &projection); err != nil {
		return
	}
	fmt.Printf("Projeção decodificada: %+v\n", projection)

	if projection.UserID == 0 || projection.Amount == 0 || projection.CategoryID == 0 || projection.Description == "" || projection.Type == "" || projection.Date.IsZero() {
		http.Error(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	if projection.IsRecurring {
		if projection.EndMonth != nil {
			endMonth, err := time.Parse("2006-01", *projection.EndMonth)
			if err != nil {
				fmt.Println("Erro ao parsear EndMonth:", err)
				http.Error(w, "Formato inválido para end_month. Use YYYY-MM.", http.StatusBadRequest)
				return
			}

			startDate := projection.Date
			for startDate.Before(endMonth.AddDate(0, 1, 0)) {
				err := createProjectionInstance(projection, startDate)
				if err != nil {
					fmt.Println("Erro ao criar projeção recorrente:", err)
					http.Error(w, "Erro ao criar projeções recorrentes", http.StatusInternalServerError)
					return
				}
				startDate = startDate.AddDate(0, 1, 0)
			}
		} else {
			http.Error(w, "Para projeções recorrentes, o campo end_month é obrigatório.", http.StatusBadRequest)
			return
		}
	} else {
		err := createProjectionInstance(projection, projection.Date)
		if err != nil {
			http.Error(w, "Erro ao inserir a projeção", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(projection)
}

func createProjectionInstance(projection models.Projection, date time.Time) error {
	query := `
        INSERT INTO projections (user_id, amount, description, category_id, type, is_recurring, end_month, date, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, false, NULL, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id
    `
	err := config.DB.QueryRow(query, projection.UserID, projection.Amount, projection.Description, projection.CategoryID, projection.Type, date).Scan(&projection.ID)
	if err != nil {
		fmt.Println("Erro ao inserir projeção:", err)
		return err
	}
	fmt.Printf("Projeção criada para o mês: %s, ID: %d\n", date.Format("2006-01"), projection.ID)
	return nil
}

func GetProjections(w http.ResponseWriter, r *http.Request) {
	projectionType := r.URL.Query().Get("type")

	query := `SELECT p.id, p.user_id, p.amount, p.description, p.type, p.date, p.created_at, p.updated_at, p.is_recurring, p.end_month, c.name 
	          FROM projections p
	          JOIN categories c ON p.category_id = c.id`

	var rows *sql.Rows
	var err error
	if projectionType != "" {
		query += " WHERE p.type = $1"
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
		var categoryName string
		err = rows.Scan(&projection.ID, &projection.UserID, &projection.Amount, &projection.Description, &projection.Type, &projection.Date, &projection.CreatedAt, &projection.UpdatedAt, &projection.IsRecurring, &projection.EndMonth, &categoryName)
		if err != nil {
			http.Error(w, "Erro ao escanear projeção", http.StatusInternalServerError)
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
		http.Error(w, "ID de projection inválido", http.StatusBadRequest)
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
		http.Error(w, "O tipo da projecao não corresponde ao tipo da categoria", http.StatusBadRequest)
		return
	}

	query := `UPDATE projections 
	          SET amount = $1, description = $2, category_id = $3, type = $4, date = $5 
	          WHERE id = $6`
	_, err = config.DB.Exec(query, projection.Amount, projection.Description, projection.CategoryID, projection.Type, projection.Date, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar a projecao", http.StatusInternalServerError)
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
