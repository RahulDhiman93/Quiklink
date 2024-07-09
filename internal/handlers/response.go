package handlers

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	resp := jsonResponse{
		OK:      false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
