package main

import (
	"Quiklink_BE/internal/config"
	"Quiklink_BE/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	gh "github.com/gorilla/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	corsHandler := gh.CORS(
		gh.AllowedOrigins([]string{"*"}),            // Allow all origins
		gh.AllowedMethods([]string{"GET", "POST"}),  // Allow only specified methods
		gh.AllowedHeaders([]string{"Content-Type"}), // Allow only specified headers
	)

	mux.Use(corsHandler)

	mux.Route("/auth", func(mux chi.Router) {
		if app.InProduction {
			mux.Use(Auth)
		}
		mux.Post("/login", handlers.Repo.LoginUser)
		mux.Post("/register", handlers.Repo.RegisterUser)
	})

	mux.Get("/", handlers.Repo.Home)
	mux.Post("/shorten", handlers.Repo.ShortenURL)
	mux.Get("/{shortKey:[a-zA-Z0-9_-]+}", handlers.Repo.Redirect)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
