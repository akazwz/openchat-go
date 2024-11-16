package xhttp

import (
	"github.com/goccy/go-json"
	"net/http"
)

func Bind(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}
