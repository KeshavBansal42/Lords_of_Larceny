package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
	"github.com/gorilla/mux"
)

func GetVillage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	value := r.Context().Value("userID")
	userIDFloat, ok := value.(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	userID := int(userIDFloat)

	_, thlevel, gold, elixir, err := repository.GetVillageByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.VillageResponseDTO{
		TownHallLevel: thlevel,
		Gold:          gold,
		Elixir:        elixir,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func CollectResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	value := r.Context().Value("userID")
	userIDFloat, ok := value.(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	userID := int(userIDFloat)

	gold, elixir, err := repository.CollectResources(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.CollectResponseDTO{
		Message: "Resources succesfully collected",
		Gold:    gold,
		Elixir:  elixir,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func ScoutVillage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetIDStr := vars["id"]
	targetUserID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		http.Error(w, "Invalid target user ID.", http.StatusBadRequest)
		return
	}

	username, thLevel, gold, elixir, buildings, err := repository.ScoutVillage(targetUserID)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
