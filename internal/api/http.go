package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) Response

func (fn HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := fn(w, r)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg(res.Message)
	}

	if err := res.Encode(w); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON response.")
	}
}
