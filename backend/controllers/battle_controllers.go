package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
)

func Matchmake(w http.ResponseWriter, r *http.Request) {
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

	villageID, err := repository.Matchmake(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.MatchmakeResponseDTO{
		Message:   "Enemy found successfully",
		VillageID: villageID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func Battle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.AttackRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid req body.", http.StatusBadRequest)
		return
	}

	value := r.Context().Value("userID")
	userIDFloat, ok := value.(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
		return
	}

	userID := int(userIDFloat)

	damage, lootedGold, lootedElixir, battleLog, err := repository.Battle(userID, req.TargetUserID, req.Drops)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.BattleResultDTO{
		PercentageDestroyed: damage,
		GoldStolen:          lootedGold,
		ElixirStolen:        lootedElixir,
		Log:                 battleLog,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
