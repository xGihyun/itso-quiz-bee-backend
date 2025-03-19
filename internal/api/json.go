package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`

	Error error `json:"-"`
}

type Status string

func (r Response) Encode(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(r.Code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		return err
	}

	return nil
}
