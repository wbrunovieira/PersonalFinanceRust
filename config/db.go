package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	retries := 10

	for retries > 0 {
		DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("Erro ao conectar no banco de dados: %v\n", err)
			retries--
			time.Sleep(3 * time.Second)

			continue
		}

		err = DB.Ping()
		if err != nil {
			fmt.Printf("Não foi possível pingar o banco de dados: %v\n", err)
			retries--
			time.Sleep(3 * time.Second)

			continue
		}

		fmt.Println("Conexão com o banco de dados foi bem-sucedida!")
		break
	}

	if retries == 0 {
		log.Fatalf("Não foi possível conectar ao banco de dados após várias tentativas. Último erro: %v\n", err)
	}

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
