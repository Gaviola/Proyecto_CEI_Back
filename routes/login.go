package routes

import (
	"encoding/json"
	"github.com/Gaviola/Proyecto_CEI_Back.git/data"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"
	"github.com/Gaviola/Proyecto_CEI_Back.git/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func LoginRoutes(r chi.Router) {
	r.Route("/login", func(r chi.Router) {
		r.Post("/user", LoginUser)     // Login con email y contraseña
		r.Post("/google", LoginGoogle) // Login con Google
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds data.Credentials
	var user data.User
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
	user, err = utils.DBExistUser(hash, creds.Username)
	if err != nil {
		if user.IsEmpty() {
			http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error de servidor", http.StatusInternalServerError)
		return
	}

	if !user.IsEmpty() {
		// Genera un token JWT
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &data.Claims{
			Username: user.Name,
			Role:     user.Role,
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
		err = json.NewEncoder(w).Encode(map[string]string{"tokenJWT": tokenString})
		if err != nil {
			http.Error(w, "No se puedo enviar el Token", http.StatusInternalServerError)
			return
		}

	}

}

func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	// Routes for the application
	http.HandleFunc("/", services.HandleMain)
	http.HandleFunc("/login-gl", services.HandleGoogleLogin)
	http.HandleFunc("/callback-gl", services.CallBackFromGoogle)
}
