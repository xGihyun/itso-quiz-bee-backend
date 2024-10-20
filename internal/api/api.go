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

	if res.Error != nil {
		log.Error().Err(res.Error).Msg(res.Message)
		http.Error(w, res.Message, res.StatusCode)
		return
	}

	w.WriteHeader(res.StatusCode)
}
