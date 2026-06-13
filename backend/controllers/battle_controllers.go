package controllers

import (
	"net/http"

	"github.com/KeshavBansal42/Lords_of_Larceny/backend/dtos"
	"github.com/KeshavBansal42/Lords_of_Larceny/backend/repository"
)

func Matchmake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	targetUserID, err := repository.Matchmake(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.MatchmakeResponseDTO{
		Message: "Enemy found successfully",
		UserID:  targetUserID,
	}

	respond(w, http.StatusOK, res)
}

func Battle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req dtos.AttackRequestDTO
	check := parseRequest(w, r, &req)
	if !check {
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	damage, lootedGold, lootedElixir, battleLog, err := repository.Battle(r.Context(), userID, req.TargetUserID, req.Drops)
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

	respond(w, http.StatusOK, res)
}

func GetBattleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(w, r)
	if err != nil {
		return
	}

	battles, err := repository.GetBattleHistory(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := dtos.GetBattleHistoryResponseDTO{
		Battles: battles,
	}

	respond(w, http.StatusOK, res)
}
