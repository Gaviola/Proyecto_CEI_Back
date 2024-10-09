package main

import (
	//"bufio"
	//"strings"

	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/configs"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/logger"

	//"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/routes"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"

	//"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	// "github.com/spf13/viper"
)

func main() {

	// Configurar CORS usando la librería rs/cors
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, // Métodos permitidos
		AllowedHeaders:   []string{"Content-Type", "Authorization", ""},                    // Cabeceras permitidas
		AllowCredentials: true,
	})

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	// Rutas de la aplicacion
	routes.LoginRoutes(r)
	routes.AdminRoutes(r)
	routes.RegisterRoutes(r)

	handler := c.Handler(r)

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
	// logger.Log.Info("Started running on http://localhost:" + viper.GetString("port")) // Log the port where the server is running
	log.Fatal(http.ListenAndServe(":8080", handler))

}
