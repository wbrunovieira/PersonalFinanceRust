package routes

import (
	"app/handlers"

	"github.com/gorilla/mux"
)

func RegisterProjectionRoutes(router *mux.Router) {
	router.HandleFunc("/projections", handlers.CreateTransaction).Methods("POST", "OPTIONS")
	router.HandleFunc("/projections", handlers.GetProjections).Methods("GET", "OPTIONS")
	router.HandleFunc("/projections/{id}", handlers.GetTransaction).Methods("GET", "OPTIONS")
	router.HandleFunc("/projections/{id}", handlers.UpdateTransaction).Methods("PUT", "OPTIONS")
	router.HandleFunc("/projections/{id}", handlers.DeleteTransaction).Methods("DELETE", "OPTIONS")
}
