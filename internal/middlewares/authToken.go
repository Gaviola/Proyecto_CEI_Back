package middlewares

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware
/*
Middleware para verificar el token de autorización en las rutas.
*/
func AuthMiddleware(next http.Handler) http.Handler {
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
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// Token válido, continuar con la solicitud
		next.ServeHTTP(w, r)
	})
}
