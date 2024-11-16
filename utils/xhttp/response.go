package xhttp

import (
	"github.com/goccy/go-json"
	"net/http"
)

func RespJson(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
