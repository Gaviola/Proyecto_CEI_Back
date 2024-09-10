package main

import (
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/configs"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/logger"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"
	"github.com/Gaviola/Proyecto_CEI_Back.git/routes"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"strings"

	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/Gaviola/Proyecto_CEI_Back.git/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)

func main() {

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	// Rutas de la aplicacion
	routes.LoginRoutes(r)

	// Initialize Viper across the application
	configs.InitializeViper()
	fmt.Println("Viper initialized...")

	// Initialize Logger across the application
	logger.InitializeZapCustomLogger()
	fmt.Println("Zap Custom Logger initialized...")

	// Initialize Oauth2 Services
	services.InitializeOAuthGoogle()
	fmt.Println("OAuth2 Services initialized...")

	//Crear llave secreta
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	secret := base64.StdEncoding.EncodeToString(key)
	err = os.Setenv("JWT_SECRET", secret)
	if err != nil {
		http.Error(nil, "Error al setear la variables de entorno", http.StatusInternalServerError)
		return
	}

	fmt.Println("Servidor escuchando en http://localhost:8080")
	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))

}

// Middleware para verificar el token de autorización en las rutas

func authMiddleware(next http.Handler) http.Handler {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		claims := &data.Claims{}
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
