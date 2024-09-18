package middlewares

import (
	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"net/http"
)

// Middleware para verificar el rol de autorización en las rutas

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar el rol del usuario del token JWT
		claims := r.Context().Value("claims").(*models.Claims)
		if claims.Role != "admin" {
			http.Error(w, "No tienes permiso para acceder a esta ruta", http.StatusForbidden)
			return
		}
		// Usuario con rol de admin, continuar con la solicitud
		next.ServeHTTP(w, r)
	})
}
