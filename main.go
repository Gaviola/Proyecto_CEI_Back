package main

import (
	"Proyecto_CEI_Back/data"
	"Proyecto_CEI_Back/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(token)
		if err != nil {
			http.Error(w, "No se pudo generar el token", http.StatusInternalServerError)
			return
		}

		// Devuelve el token al cliente
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}

}

func admin(w http.ResponseWriter, r *http.Request) {
	//logica de admin
}

func main() {
	fmt.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	http.HandleFunc("/", LoginHandler)

}
