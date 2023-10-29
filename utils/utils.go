package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, status int, err Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
