package routes

import (
	"app/handlers"

	"github.com/gorilla/mux"
)

func RegisterCategoryRoutes(router *mux.Router) {
	router.HandleFunc("/categories", handlers.GetCategories).Methods("GET", "OPTIONS")
	router.HandleFunc("/categories", handlers.CreateCategory).Methods("POST", "OPTIONS")
	router.HandleFunc("/categories/{id}", handlers.GetCategory).Methods("GET", "OPTIONS")
	router.HandleFunc("/categories/{id}", handlers.UpdateCategory).Methods("PUT", "OPTIONS")
	router.HandleFunc("/categories/{id}", handlers.DeleteCategory).Methods("DELETE", "OPTIONS")
}
