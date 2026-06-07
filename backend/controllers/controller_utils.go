package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func respond(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func getUserID(w http.ResponseWriter, r *http.Request) (int, error) {
	value := r.Context().Value("userID")
	userIDFloat, ok := value.(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return 0, errors.New("unauthorized")
	}
	return int(userIDFloat), nil
}

func parseRequest(w http.ResponseWriter, r *http.Request, req any) bool {
	err := json.NewDecoder(r.Body).Decode(req) // Notice req doesn't have an & here
	if err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return false
	}
	return true
}
