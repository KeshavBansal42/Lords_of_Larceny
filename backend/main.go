package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB()

	r := mux.NewRouter()
	routes.InitRoutes(r)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	addrString := os.Getenv("SERVER_URI")

	http.ListenAndServe(addrString, r)
}
