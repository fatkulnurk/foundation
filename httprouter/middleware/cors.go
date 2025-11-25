package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

type CORSOptions struct {
	AllowedOrigins   []string // ["*"] untuk semua origin
	AllowedMethods   []string // contoh: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	AllowedHeaders   []string // contoh: []string{"Content-Type", "Authorization"}
	ExposedHeaders   []string // optional
	AllowCredentials bool
	MaxAge           int // detik, contoh: 600
}

// CORS mengembalikan middleware CORS net/http
func CORS(opts CORSOptions) func(http.Handler) http.Handler {
	allowedMethods := joinOrDefault(opts.AllowedMethods,
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	)
	allowedHeaders := joinOrDefault(opts.AllowedHeaders,
		[]string{"Content-Type", "Authorization"},
	)

	allowedMethodsStr := strings.Join(allowedMethods, ", ")
	allowedHeadersStr := strings.Join(allowedHeaders, ", ")
	exposedHeadersStr := strings.Join(opts.ExposedHeaders, ", ")
	anyOrigin := contains(opts.AllowedOrigins, "*")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// bukan request CORS
				next.ServeHTTP(w, r)
				return
			}

			// Tentukan origin yang di-allow
			if anyOrigin {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if contains(opts.AllowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			} else {
				// origin tidak diizinkan -> lanjut tanpa header CORS
				next.ServeHTTP(w, r)
				return
			}

			if opts.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if exposedHeadersStr != "" {
				w.Header().Set("Access-Control-Expose-Headers", exposedHeadersStr)
			}

			// Preflight request
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", allowedMethodsStr)
				w.Header().Set("Access-Control-Allow-Headers", allowedHeadersStr)
				if opts.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(opts.MaxAge))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func joinOrDefault(got []string, def []string) []string {
	if len(got) == 0 {
		return def
	}
	return got
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
