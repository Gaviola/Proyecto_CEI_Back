package main

import (
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/configs"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/logger"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/routes"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"crypto/rand"
	"encoding/base64"
	"fmt"
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
	routes.RegisterRoutes(r)
	routes.AdminRoutes(r)

	// Initialize Viper across the application
	configs.InitializeViper()
	fmt.Println("Viper initialized...")

	// Initialize Logger across the application
	logger.InitializeZapCustomLogger()
	fmt.Println("Zap Custom Logger initialized...")

	// Initialize Oauth2 Services
	services.InitializeOAuthGoogle()
	fmt.Println("OAuth2 Services initialized...")

	var opcion int
	fmt.Println("DEBUG MENU")
	fmt.Println("[1] \tSetear variables de entorno")
	fmt.Println("[2] \tGoogle Login")
	fmt.Println("[3] \tPrueba de ItemTypes")
	fmt.Println("[4] \tPrueba de Items")
	fmt.Println("[OTRO] \tSalir")

	fmt.Print("> Ingrese una opción: ")
	fmt.Scan(&opcion)
	fmt.Println()

	switch opcion {
	case 1:
		fmt.Println("Has seleccionado la opción 1")
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

	case 2:
		// Routes for the application
		http.HandleFunc("/", services.HandleMain)
		http.HandleFunc("/login-gl", services.HandleGoogleLogin)
		http.HandleFunc("/callback-gl", services.CallBackFromGoogle)

		log.Fatal(http.ListenAndServe(":8080", nil))

		fmt.Println("Servidor escuchando en http://localhost:8080")
		logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
		log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))

	case 3:
		itemTypes := repositories.DBShowItemTypes()
		for _, itemType := range itemTypes {
			fmt.Println(itemType)
		}

	case 4:
		items := repositories.DBShowItems()
		for _, item := range items {
			fmt.Println(item)
		}
	default:
		fmt.Println("Saliendo...")
	}

}
