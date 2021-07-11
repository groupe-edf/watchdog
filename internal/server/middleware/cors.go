package middleware

import (
	"net/http"
	"strings"
)

// CORS middleware
type CORS struct {
}

func (middleware *CORS) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := []string{
			"Accept",
			"Accept-Ranges",
			"Authorization",
			"Content-Range",
			"Content-Type",
			"Origin",
			"X-Requested-With",
			"X-Token",
		}
		if r.Method == http.MethodOptions {
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Add("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(headers, ", "))
		next.ServeHTTP(w, r)
	})
}
