package middleware

import (
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

func Rate(log *slog.Logger, next http.HandlerFunc, rps int) http.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), 1)
	return func(w http.ResponseWriter, r *http.Request) {
		if err := limiter.Wait(r.Context()); err != nil {
			log.Error("Rate limit exceeded", "error", err)
			http.Error(w, "server is going down", http.StatusServiceUnavailable)
			return
		}
		next.ServeHTTP(w, r)
	}
}
