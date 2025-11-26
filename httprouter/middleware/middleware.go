package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// middleware pakai net/http murni
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println(">>", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Println("<<", r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				log.Printf("[PANIC] %v\n%s", err, debug.Stack())

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
