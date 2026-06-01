package routes

import (
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/controllers"
	"github.com/gorilla/mux"
)

func InitRoutes(router *mux.Router) {
	router.HandleFunc("/register", controllers.Register).Methods("POST")
}
