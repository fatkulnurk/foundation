package middleware

import "net/http"

// RequireAPIKey only for example
func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") == "" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
