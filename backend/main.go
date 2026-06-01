package main

import (
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/db"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	r := mux.NewRouter()

}
