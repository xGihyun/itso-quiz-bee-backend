package api

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) Response

func (fn HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := fn(w, r)

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		res.Status = Success
	} else if res.StatusCode >= 400 && res.StatusCode <= 499 {
		res.Status = Fail
	} else {
		res.Status = Error
	}

	if res.Error != nil {
		if res.Message == "" {
			res.Message = res.Error.Error()
		}

		log.Error().Err(res.Error).Msg(res.Message)
	}

	if err := res.Encode(w); err != nil {
		log.Error().Err(err).Msg("Failed to encode JSON response.")
	}
}
