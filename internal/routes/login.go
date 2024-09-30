package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/services"
	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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
		expirationTime := time.Now().Add(30 * time.Minute)
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

		// If the user does not exist, send a message to the client to register the user and send the mail
		if user.IsEmpty() {
			err = json.NewEncoder(w).Encode(map[string]string{"message": "User does not exist", "email": googleUser.Email})
			if err != nil {
				http.Error(w, "Error de servidor", http.StatusInternalServerError)
			}
		} else {
			print("User exists\n")
		}
	})

}

func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	//Recibo el mail del usuario y lo guardo en un usuario temporal
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// busco al usuario en base al mail
	var foundUser models.User
	foundUser, err = repositories.DBGetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "Error de servidor", http.StatusInternalServerError)
		return
	}
	if foundUser.IsEmpty() {
		http.Error(w, "Usuario no registrado", http.StatusUnauthorized)
		return
	} else {
		//Genero un token para el usuario
		var token string
		token, err := generateResetToken(foundUser)
		if err != nil {
			http.Error(w, "Error al generar el token", http.StatusInternalServerError)
			return
		}
		err = sendPasswordResetEmail(foundUser.Email, token)
		if err != nil {
			http.Error(w, "Error al enviar el mail", http.StatusInternalServerError)
			return
		}
	}

}

func generateResetToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &models.Claims{
		Username: strconv.Itoa(user.ID),
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// utilizar llave secreta para firmar el token
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("no se ha definido una llave secreta")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("no se pudo generar el token")
	}
	return tokenString, nil
}

func sendPasswordResetEmail(toEmail string, token string) error {
	resetURL := "https://CEI-Prestamos/restablecer-contraseña?token=" + token
	body := fmt.Sprintf("Click here to reset your password: %s", resetURL)

	// Configura el servidor SMTP
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@CEI-Prestamos.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Restablecimiento de contraseña")
	m.SetBody("text/plain", body)

	// Conexión con el servidor SMTP. (Deberiamos utilizar un mail del cei o algo asi. Guardar la contraseña en una variable de entorno)
	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-password")

	// Enviar el correo.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	// Obtener el token de la URL.
	tokenString := r.URL.Query().Get("token")

	// Parsear y verificar el token.
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// El token es válido, obtenemos el ID del usuario desde los claims.
	userID, err := strconv.Atoi(claims.Username)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Ahora puedes permitir al usuario cambiar la contraseña.
	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hashear la nueva contraseña y actualizarla en la base de datos.
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// Actualiza la contraseña en la base de datos.
	user := models.User{
		Hash: hash,
	}
	err = repositories.DBUpdateUser(userID, user)
	if err != nil {
		http.Error(w, "Could not update password", http.StatusInternalServerError)
		return
	}

	// Responder éxito.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been reset successfully",
	})
	if err != nil {
		return
	}
}
