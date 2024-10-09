package middlewares

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/dgrijalva/jwt-go"
)

// AdminMiddleware
/*
Middleware para verificar el rol de autorización en las rutas.
*/
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtKey []byte
		jwtKey, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_SECRET"))
		if err != nil {
			http.Error(w, "Error al decodificar la llave secreta", http.StatusInternalServerError)
			return
		}

		// Obtener el token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Falta el header Authorization", http.StatusUnauthorized)
			return
		}

		// Verificar que el formato sea "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato de autorización inválido", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Parsear y validar el token
		claims := &models.Claims{}
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}
		if !(claims.Role == "admin") {
			http.Error(w, "No tienes permisos de administrador", http.StatusUnauthorized)
			return
		}

		// Usuario con rol de admin, continuar con la solicitud
		next.ServeHTTP(w, r)
	})
}
