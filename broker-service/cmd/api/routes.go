package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

func (app *Config) routes() (handler http.Handler) {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// health check
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
			log.Println()
			next.ServeHTTP(w, r)
		})
	})
	//
	mux.Post("/", app.Broker)
	//mux.Post("/log-grpc", app.logItemViaGPRC)
	//mux.Post("/handle", app.HandleSubmission)

	handler = mux
	return
}
