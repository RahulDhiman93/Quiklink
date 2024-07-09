package main

import (
	"Quiklink_BE/internal/helpers"
	gh "github.com/gorilla/handlers"
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf Use for CSRF Token
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad Use for loading Session
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CorsHandler(next http.Handler) http.Handler {
	cors := gh.CORS(
		gh.AllowedOrigins([]string{"*"}),            // Allow all origins
		gh.AllowedMethods([]string{"GET", "POST"}),  // Allow only specified methods
		gh.AllowedHeaders([]string{"Content-Type"}), // Allow only specified headers
	)
	return cors(next)
}
