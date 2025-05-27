package middleware

import (
	"log/slog"

	"net/http"
	"strings"
)

type TokenVerifier interface {
	Verify(token string) error
}

func Auth(log *slog.Logger, f http.HandlerFunc, verifier TokenVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		log.Debug("Got token", "token", parts)

		if len(parts) != 2 || parts[0] != "Token" {
			http.Error(w, "bad authorization header", http.StatusUnauthorized)
			return
		}
		if err := verifier.Verify(parts[1]); err != nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}
		f.ServeHTTP(w, r)
	}
}
