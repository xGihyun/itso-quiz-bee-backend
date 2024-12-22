package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func RequestLogger(next http.Handler) http.Handler {
	h := hlog.NewHandler(log.Logger)

	access := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {

			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Int("status_code", status).
				Str("user_agent", r.UserAgent()).
				Dur("elapsed_ms", duration).
				Msg("Incoming request.")
		},
	)

	userAgent := hlog.UserAgentHandler("user_agent")

	return h(access(userAgent(next)))
}
