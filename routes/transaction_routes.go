package routes

import (
	"app/handlers"

	"github.com/gorilla/mux"
)

func RegisterTransactionRoutes(router *mux.Router) {
	router.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST")
	router.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/{id}", handlers.GetTransaction).Methods("GET")
	router.HandleFunc("/transactions/{id}", handlers.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", handlers.DeleteTransaction).Methods("DELETE")
}
