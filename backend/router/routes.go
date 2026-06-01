package router

import (
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/controllers"
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", controllers.Register).Methods("POST")

	return router
}
