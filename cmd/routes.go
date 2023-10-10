package main

import (
	"Quiklink_BE/internal/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/shorten", handlers.Repo.ShortenURL)
	mux.Get("/{shortKey:[a-zA-Z0-9]+}", handlers.Repo.Redirect)

	return mux
}
