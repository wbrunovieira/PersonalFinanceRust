package routes

import (
	"app/handlers"

	"github.com/gorilla/mux"
)

func RegisterTransactionRoutes(router *mux.Router) {
	router.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST", "OPTIONS")
	router.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET", "OPTIONS")
	router.HandleFunc("/transactions/{id}", handlers.GetTransaction).Methods("GET", "OPTIONS")
	router.HandleFunc("/transactions/{id}", handlers.UpdateTransaction).Methods("PUT", "OPTIONS")
	router.HandleFunc("/transactions/{id}", handlers.DeleteTransaction).Methods("DELETE", "OPTIONS")
}
