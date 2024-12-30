package routes

import (
	"app/handlers"

	"github.com/gorilla/mux"
)

func RegisterProjectionRoutes(router *mux.Router) {
	router.HandleFunc("/projections", handlers.CreateProjection).Methods("POST", "OPTIONS")
	router.HandleFunc("/projections", handlers.GetProjections).Methods("GET", "OPTIONS")
	router.HandleFunc("/projections/{id}", handlers.GetProjection).Methods("GET", "OPTIONS")
	router.HandleFunc("/projections/{id}", handlers.UpdateProjection).Methods("PUT", "OPTIONS")
	router.HandleFunc("/projections/{id}", handlers.DeleteProjection).Methods("DELETE", "OPTIONS")
}
