package routes

import (
	"encoding/json"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/middlewares"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func AdminRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)  // middleware de verificacion de token
		r.Use(middlewares.AdminMiddleware) // middleware de verificacion de admin
		r.Post("/", CreateUser)            // Crear un usuario
		r.Delete("/{userID}", DeleteUser)  // Eliminar un usuario
		r.Patch("/{userID}", UpdateUser)   // Actualizar un usuario
		r.Get("/", GetUsers)               // Obtener todos los usuarios
		r.Get("/{userID}", GetUser)        // Obtener un usuario
		r.Put("/{userID}", VerifyUser)     // Verificar un usuario
		//TODO manejo de items
		//TODO manejo de prestamos
	})
}

// CreateUser
/*
Recibe los datos de un nuevo usuario, verifica que los campos sean correctos y los guarda en la base de datos
*/
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Crea un nuevo usuario a partir de los datos recibidos
	var user models.User
	// tal vez explota. falta probar
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Maneja el error si los datos no son correctos
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	// Verificar que el usuario no exista
	var existUser bool
	existUser, err = repositories.DBCheckUser(user.Email)
	if err != nil {
		http.Error(w, "Error checking user", http.StatusInternalServerError)
		return
	}
	if !existUser {
		// Guarda el nuevo usuario en la base de datos
		err = repositories.DBSaveUser(user)
		if err != nil {
			return
		}
		// Responde con el nuevo usuario creado
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			return
		}
	} else {
		http.Error(w, "User already exists", http.StatusBadRequest)
	}

}

// DeleteUser
/*
Recibe el id de un usuario y lo elimina de la base de datos
*/
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Captura el valor del parámetro userID de la URL
	userID := chi.URLParam(r, "userID")
	// Convierte el userID a un número entero, si es necesario
	id, err := strconv.ParseInt(userID, 10, 0)
	if err != nil {
		// Maneja el error si el userID no es un número válido
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Elimina el usuario de la base de datos
	err = repositories.DBDeleteUser(int(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
}

// UpdateUser
/*
Recibe los datos de un usuario y los actualiza en la base de datos
*/

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Captura el valor del parámetro userID de la URL
	userID := chi.URLParam(r, "userID")
	// Convierte el userID a un número entero, si es necesario
	id, err := strconv.ParseInt(userID, 10, 0)
	if err != nil {
		// Maneja el error si el userID no es un número válido
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// busca el usuario en la base de datos
	user, err := repositories.DBGetUserByID(int(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Actualiza los datos del usuario con los datos recibidos
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Maneja el error si los datos no son correctos
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

}

// GetUsers
/*
Obtiene todos los usuarios de la base de datos
*/
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repositories.DBGetAllUsers()
	if err != nil {
		return
	}
	// Responde con los usuarios encontrados
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		return
	}
}

// GetUser
/*
Obtiene un usuario de la base de datos
*/
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Captura el valor del parámetro userID de la URL
	userID := chi.URLParam(r, "userID")
	// Convierte el userID a un número entero, si es necesario
	id, err := strconv.ParseInt(userID, 10, 0)
	if err != nil {
		// Maneja el error si el userID no es un número válido
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := repositories.DBGetUserByID(int(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Responde con el usuario encontrado
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// VerifyUser
/*
Recibe el id de un usuario y lo verifica
*/
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	// Captura el valor del parámetro userID de la URL
	userID := chi.URLParam(r, "userID")
	// Convierte el userID a un número entero, si es necesario
	id, err := strconv.ParseInt(userID, 10, 0)
	if err != nil {
		// Maneja el error si el userID no es un número válido
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
}
