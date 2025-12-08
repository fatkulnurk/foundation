package middleware

import "net/http"

// RequireAPIKey is a middleware that checks for X-API-Key header
func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}
		// In production, you would validate the API key against a database or cache
		// For this example, we just check if it exists
		next.ServeHTTP(w, r)
	})
}
