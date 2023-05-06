package router

import (
	handlers "backend/internal/handler"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	// Initialize the HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/users/{user_id:[0-9]+}", handlers.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/produce/{user_id:[0-9]+}", handlers.ProduceUserEventHandler).Methods("POST")

	return router
}
