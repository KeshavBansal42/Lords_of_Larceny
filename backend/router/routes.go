package router

import (
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/controllers"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/middleware"
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/game/configs", controllers.GetGameConfigs).Methods("GET")
	router.HandleFunc("/user/delete", controllers.DeleteAccount).Methods("DELETE")

	villageRouter := router.PathPrefix("/village").Subrouter()
	villageRouter.Use(middleware.RequireAuth)

	villageRouter.HandleFunc("/", controllers.GetVillage).Methods("GET")
	villageRouter.HandleFunc("/{id}/scout", controllers.ScoutVillage).Methods("GET")

	villageRouter.HandleFunc("/buildings", controllers.GetAllVillageBuildings).Methods("GET")
	villageRouter.HandleFunc("/buildings/build", controllers.AddBuilding).Methods("POST")
	villageRouter.HandleFunc("/buildings/upgrade", controllers.UpgradeBuilding).Methods("PUT")
	villageRouter.HandleFunc("/buildings/move", controllers.MoveBuilding).Methods("PUT")

	villageRouter.HandleFunc("/troops", controllers.GetAllVillageTroops).Methods("GET")
	villageRouter.HandleFunc("/troops/train", controllers.TrainTroops).Methods("PUT")

	villageRouter.HandleFunc("/collect/gold", controllers.CollectGold).Methods("PUT")
	villageRouter.HandleFunc("/collect/elixir", controllers.CollectElixir).Methods("PUT")

	battleRouter := router.PathPrefix("/battle").Subrouter()
	battleRouter.Use(middleware.RequireAuth)

	battleRouter.HandleFunc("/matchmake", controllers.Matchmake).Methods("GET")
	battleRouter.HandleFunc("/attack", controllers.Battle).Methods("POST")
	battleRouter.HandleFunc("/history", controllers.GetBattleHistory).Methods("GET")

	return router
}
