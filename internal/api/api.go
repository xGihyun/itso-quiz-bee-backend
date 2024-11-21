package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) Response

func (fn HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	res := fn(w, r)

	log.Debug().Msg("I AM HERE")
	if res.Error != nil {
		if res.Message == "" {
			res.Message = res.Error.Error()
		}

		if res.Status == "" {
			res.Status = Error
		}

		log.Error().Err(res.Error).Msg(res.Message)
	} else {
		// res.Status = Success
		// res.StatusCode = 200
	}

	if err := res.Encode(w); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON error response.")
		return
	}
}
