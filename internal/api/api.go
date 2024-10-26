package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Error      error  `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type HTTPHandler func(w http.ResponseWriter, r *http.Request) Response

func (fn HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := fn(w, r)

	if res.StatusCode == 0 {
		res.StatusCode = 200
	}

	// if err := WriteJSON(w, res); err != nil {
	// 	log.Error().Err(err).Msg("Failed to encode JSON error response.")
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	if res.Error != nil {
		if res.Message == "" {
			res.Message = res.Error.Error()
		}

		log.Error().Err(res.Error).Msg(res.Message)
		http.Error(w, res.Message, res.StatusCode)
		return
	}

	w.WriteHeader(res.StatusCode)
}
