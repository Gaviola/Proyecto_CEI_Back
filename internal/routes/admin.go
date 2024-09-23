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

		r.Post("/createUser", CreateUser)  // Crear un usuario
		r.Put("/deleteUser", DeleteUser)   // Eliminar un usuario
		r.Put("/updateUser", UpdateUser)   // Actualizar un usuario
		r.Get("/getUsers", GetUsers)       // Obtener todos los usuarios
		r.Get("/getUser", GetUser)         // Obtener un usuario
		r.Put("/verifyUser", VerifyUser)   // Verificar un usuario

		r.Post("/createItemType", createItemType) 	// Crear un tipo de item
		r.Put("/deleteItemType", DeleteItemType)	// Eliminar un tipo de item
		r.Put("/updateItemType", UpdateItemType)	// Actualizar un tipo de item
		r.Get("/getItemTypes", GetItemTypes)		// Obtener todos los tipos de item
		r.Get("/getItemType", GetItemType)			// Obtener un tipo de item

		r.Post("/createItem", CreateItem)	// Crear un item
		r.Put("/deleteItem", DeleteItem)	// Eliminar un item
		r.Put("/updateItem", UpdateItem)	// Actualizar un item
		r.Get("/getItems", GetItems)		// Obtener todos los items
		r.Get("/getItem", GetItem)			// Obtener un item

		//TODO manejo de prestamos
	})
}

// CreateUser
/*
Recibe los datos de un nuevo usuario,
verifica que los campos sean correctos
y los guarda en la base de datos
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

// CreateItemType
/*
Recibe los datos de un nuevo tipo de item,
verifica que los campos sean correctos
y los guarda en la base de datos
*/
func createItemType(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// DeleteItemType
/*
Recibe el id de un tipo de item y lo elimina de la base de datos
*/
func DeleteItemType(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// UpdateItemType
/*
Recibe los datos de un tipo de item y los actualiza en la base de datos
*/
func UpdateItemType(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// GetItemTypes
/*
Obtiene todos los tipos de item de la base de datos
*/
func GetItemTypes(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// GetItemType
/*
Obtiene un tipo de item de la base de datos
*/
func GetItemType(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// CreateItem
/*
Recibe los datos de un nuevo item,
verifica que los campos sean correctos
y los guarda en la base de datos
*/
func CreateItem(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// DeleteItem
/*
Recibe el id de un item y lo elimina de la base de datos
*/
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// UpdateItem
/*
Recibe los datos de un item y los actualiza en la base de datos
*/
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// GetItems
/*
Obtiene todos los items de la base de datos
*/
func GetItems(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// GetItem
/*
Obtiene un item de la base de datos
*/
func GetItem(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}
