package middlewares

import (
	"net/http"

	"github.com/NischithB/chirpy/config"
)

func FileServerHitCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)

		config.Config.FileServerHits++
	})
}
