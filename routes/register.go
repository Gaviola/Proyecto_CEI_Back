package routes

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/register", func(r chi.Router) {
		r.Post("/user", LoginUser)     // registro manual
		r.Post("/google", LoginGoogle) // registo con Google (requiere de verificaciones extras)
	})
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

func RegisterGoogle(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}
