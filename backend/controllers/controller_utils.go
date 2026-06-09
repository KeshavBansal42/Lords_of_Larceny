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

func getUserID(w http.ResponseWriter, r *http.Request) (string, error) {
	value := r.Context().Value("userID")
	userIDStr, ok := value.(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return "", errors.New("unauthorized")
	}
	return userIDStr, nil
}

func parseRequest(w http.ResponseWriter, r *http.Request, req any) bool {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return false
	}
	return true
}
