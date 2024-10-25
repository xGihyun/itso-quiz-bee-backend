package api

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}

	return nil
}
