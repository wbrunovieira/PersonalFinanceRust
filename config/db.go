package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v\n", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Não foi possível pingar o banco de dados: %v\n", err)
	}

	fmt.Println("Conexão com o banco de dados foi bem-sucedida!")

	runMigrations()
}

func runMigrations() {
	fmt.Println("Iniciando as migrações do banco de dados...")

	fmt.Printf("Path: %s\n", "/app/db/migrations")
	fmt.Printf("DATABASE_URL: %s\n", os.Getenv("DATABASE_URL"))

	cmd := exec.Command("migrate", "-path", "/app/db/migrations", "-database", os.Getenv("DATABASE_URL"), "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Erro ao executar migrações: %v\nSaída: %s\n", err, string(output))
	}
	fmt.Println("Migrações do banco de dados concluídas com sucesso.")
}
