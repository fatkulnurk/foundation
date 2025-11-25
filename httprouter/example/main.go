package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatkulnurk/foundation/httprouter"
	"github.com/fatkulnurk/foundation/httprouter/middleware"
)

func main() {
	r := httprouter.New()
	// global middleware
	r.Use(middleware.Logging)

	// global recovery
	r.Use(middleware.RecoverMiddleware)

	// cors
	r.Use(middleware.CORS(middleware.CORSOptions{
		AllowedOrigins:   []string{"*"}, // atau []string{"https://app.zeedsharia.com"}
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           600,
	}))

	// /static
	r.Static("/static", "./public/static")

	// rate limit
	r.Use(middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
		Requests: 10,
		Window:   time.Minute,
	}))

	// route tanpa group
	r.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		httprouter.WriteHTML(w, http.StatusOK, "pong")
	})

	r.GET("/danger", func(w http.ResponseWriter, r *http.Request) {
		panic("boom!")
	}, middleware.RecoverMiddleware)

	// group /api
	r.Group("/api", func(api httprouter.HttpRouter) {

		// middleware khusus group
		api.Use(middleware.RequireAPIKey)

		// serve /api/assets/* (dengan middleware di atas)
		api.Static("/assets", "./public/app-assets")

		// GET /api/users/{id}
		api.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id") // native Go 1.22+
			fmt.Fprintf(w, "User ID = %s\n", id)
		})

		// GET /api/admin/stats (butuh X-API-Key + X-Role=admin)
		api.Group("/admin", func(admin httprouter.HttpRouter) {
			admin.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Header.Get("X-Role") != "admin" {
						http.Error(w, "forbidden", http.StatusForbidden)
						return
					}
					next.ServeHTTP(w, r)
				})
			})

			admin.GET("/stats", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "GET /api/v1/admin/stats")
			})
		})
	})

	log.Println("listen :18080")
	log.Fatal(http.ListenAndServe(":18080", r))
}
