package handlers

import (
	"encoding/json"
	"net/http"
)

func internalServerError(w http.ResponseWriter, err error) {
	resp := jsonResponse{
		OK:      false,
		Message: err.Error(),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
