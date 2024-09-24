package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"
	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginRoutes
/*
LoginRoutes define las rutas para el login de usuarios.
*/
func LoginRoutes(r chi.Router) {
	r.Route("/login", func(r chi.Router) {
		r.Post("/user", LoginUser)     // Login con email y contraseña
		r.Post("/google", LoginGoogle) // Login con Google
	})
}

// LoginUser
/*
LoginUser permite a un usuario autenticarse con email y contraseña.
*/
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	var user models.User
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
	user, err = repositories.DBExistUser(hash, creds.Username)
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
		claims := &models.Claims{
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

// LoginGoogle
/*
LoginGoogle permite a un usuario autenticarse con Google.
*/
func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	// Routes for the application
	http.HandleFunc("/", services.HandleMain)
	http.HandleFunc("/login-gl", services.HandleGoogleLogin)
	http.HandleFunc("/callback-gl", func(w http.ResponseWriter, r *http.Request) {
		response, err := services.CallBackFromGoogle(w, r)
		if err != nil {
			http.Error(w, "Error en la autenticación con Google", http.StatusInternalServerError)
			return
		}

		// Parse the response from Google
		var googleUser models.GoogleUser
		err = json.Unmarshal(response, &googleUser)
		if err != nil {
			http.Error(w, "Error en la autenticación con Google", http.StatusInternalServerError)
			return
		}

		// Check if the user exists in the database
		var user models.User
		user, err = repositories.DBGetUserByEmail(googleUser.Email)
	
		if err != nil {
			http.Error(w, "Error de servidor", http.StatusInternalServerError)
			return
		}

		// If the user does not exist, create a new user
		if user.IsEmpty() {
			user = models.User{
				Name:       "",
				Lastname:   "",
				StudentId:  0,
				Email:      googleUser.Email,
				Phone:      0,
				Role:       "user",
				Dni:        0,
				CreatorId:  0,
				School:     "",
				IsVerified: false,
				Hash:       []byte(""), // Empty hash
			}

			err = repositories.DBSaveUser(user)
			if err != nil {
				http.Error(w, "Error de servidor", http.StatusInternalServerError)
				return
			}
		} else {
			print("User exists\n")
		}
	})

}
