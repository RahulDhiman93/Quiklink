package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func internalServerError(w http.ResponseWriter, err error) {
	resp := errorResponse{
		OK:      false,
		Message: err.Error(),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
