package main

import (
	"github.com/Gaviola/Proyecto_CEI_Back.git/data"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/configs"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/logger"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"
	"github.com/Gaviola/Proyecto_CEI_Back.git/utils"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/spf13/viper"
)

func main() {

	http.HandleFunc("/login-usr", LoginUser)

	// Initialize Viper across the application
	configs.InitializeViper()
	fmt.Println("Viper initialized...")

	// Initialize Logger across the application
	logger.InitializeZapCustomLogger()
	fmt.Println("Zap Custom Logger initialized...")

	// Initialize Oauth2 Services
	services.InitializeOAuthGoogle()
	fmt.Println("OAuth2 Services initialized...")

	// Routes for the application
	http.HandleFunc("/", services.HandleMain)
	http.HandleFunc("/login-gl", services.HandleGoogleLogin)
	http.HandleFunc("/callback-gl", services.CallBackFromGoogle)

	//Crear llave secreta
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}
	secret := base64.StdEncoding.EncodeToString(key)
	os.Setenv("JWT_SECRET", secret)
	print(secret)

	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("Servidor escuchando en http://localhost:8080")

	logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds data.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Busca el usuario en la base de datos segun el hash de la contraseña y username
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error de servidor", http.StatusInternalServerError)
		return
	}
	var isValid bool
	//Existe la posibilidad que haya colisiones??
	isValid, err = utils.DBExistUser(hash, creds.Username)
	if err != nil {
		if !isValid {
			http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error de servidor", http.StatusInternalServerError)
		return
	}

	if isValid {
		// Genera un token JWT
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &data.Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// utilizar llave secreta para firmar el token
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			http.Error(w, "No se ha definido una llave secreta", http.StatusInternalServerError)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			http.Error(w, "No se pudo generar el token", http.StatusInternalServerError)
			return
		}

		// Enviar el token en la respuesta
		json.NewEncoder(w).Encode(map[string]string{"tokenJWT": tokenString})

	}

}

func LoginTokenJWS(w http.ResponseWriter, r *http.Request) {
	//logica de login con token JWS.

	//Obtener el token JWS
	var tokenJWS string
	err := json.NewDecoder(r.Body).Decode(&tokenJWS)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// TODO terminar el coso

}
