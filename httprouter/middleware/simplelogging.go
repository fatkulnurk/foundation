package middleware

import (
	"log"
	"net/http"
	"time"
)

func SimpleLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println(">>", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Println("<<", r.Method, r.URL.Path, time.Since(start))
	})
}
