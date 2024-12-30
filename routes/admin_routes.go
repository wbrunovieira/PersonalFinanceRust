package routes

import (
	"app/handlers"
	"log"

	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(r *mux.Router) {
	log.Println("Registrando rota /admin/reset-db")
	r.HandleFunc("/admin/reset-db", handlers.ResetDatabase).Methods("POST")
}
