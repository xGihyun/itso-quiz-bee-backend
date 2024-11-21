package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Error      error  `json:"-"`
	Status     Status `json:"status"`
	Data       any    `json:"data"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type Status string

const (
	Success Status = "success"
	Error   Status = "error"
	Fail    Status = "fail"
)

func (r Response) Encode(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	if r.Status != Success {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	w.WriteHeader(r.StatusCode)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		return err
	}

	return nil
}
