package handlers

import (
	"app/config"
	"net/http"
	"os"
)

func ResetDatabase(w http.ResponseWriter, r *http.Request) {
	_, err := config.DB.Exec(`
        DROP SCHEMA public CASCADE;
        CREATE SCHEMA public;
    `)
	if err != nil {
		http.Error(w, "Erro ao resetar o banco de dados", http.StatusInternalServerError)
		return
	}

	migrationFiles := []string{
		"./db/migrations/001_create_users_table.up.sql",
		"./db/migrations/002_create_transactions_table.up.sql",
		"./db/migrations/003_create_categories_table.up.sql",
		"./db/migrations/004_create_projections_table.up.sql",
	}

	for _, file := range migrationFiles {
		sql, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "Erro ao ler arquivo de migração: "+file, http.StatusInternalServerError)
			return
		}

		_, err = config.DB.Exec(string(sql))
		if err != nil {
			http.Error(w, "Erro ao executar migração: "+file, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Banco de dados resetado com sucesso!"))
}
