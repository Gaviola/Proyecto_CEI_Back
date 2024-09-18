package routes

import (
	"encoding/json"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/register", func(r chi.Router) {
		r.Post("/user", LoginUser) // registro manual
	})
}

// RegisterUser
/*
Recibe los datos de un nuevo usuario, verifica que los campos sean correctos y los guarda en la base de datos
*/
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//Recibo los datos del usuario desde el frontend y los guardo en la base de datos
	newUser := models.User{}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Verifico que no se encuentre ya registrado
	existingUser, err := repositories.DBGetUserByEmail(newUser.Email)
	if err != nil {
		http.Error(w, "Error checking email", http.StatusInternalServerError)
		return
	}
	if !existingUser.IsEmpty() {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	newUser.IsVerified = false
	err = repositories.DBSaveUser(newUser)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "User registered, pending approval"})
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
