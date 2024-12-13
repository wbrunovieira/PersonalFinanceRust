package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"app/config"
	"app/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.InitDB()
	defer config.DB.Close()

	router := mux.NewRouter()

	routes.RegisterUserRoutes(router)
	routes.RegisterTransactionRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
