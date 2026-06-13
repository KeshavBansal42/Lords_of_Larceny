package controllers

import (
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
	"github.com/gorilla/mux"
)

func GetVillage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	err = repository.SyncBuildings(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, thlevel, gold, elixir, err := repository.GetVillageByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.VillageResponseDTO{
		TownHallLevel: thlevel,
		Gold:          gold,
		Elixir:        elixir,
	}

	respond(w, http.StatusOK, res)
}

func CollectGold(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	err = repository.SyncBuildings(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gold, err := repository.CollectGold(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.CollectGoldResponseDTO{
		Message: "Gold succesfully collected",
		Gold:    gold,
	}

	respond(w, http.StatusOK, res)
}

func CollectElixir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	err = repository.SyncBuildings(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	elixir, err := repository.CollectElixir(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.CollectElixirResponseDTO{
		Message: "Elixir succesfully collected",
		Elixir:  elixir,
	}

	respond(w, http.StatusOK, res)
}

func ScoutVillage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	err := repository.SyncBuildings(r.Context(), targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, thLevel, gold, elixir, buildings, err := repository.ScoutVillage(r.Context(), targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := dtos.ScoutVillageResponseDTO{
		Username:      username,
		TownHallLevel: thLevel,
		Gold:          gold,
		Elixir:        elixir,
		Buildings:     buildings,
	}

	respond(w, http.StatusOK, res)
}
