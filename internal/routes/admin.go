package routes

import (
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func AdminRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)  // middleware de verificacion de token
		r.Use(middlewares.AdminMiddleware) // middleware de verificacion de admin
		r.Post("/createUser", CreateUser)  // Crear un usuario
		r.Put("/deleteUser", DeleteUser)   // Eliminar un usuario
		r.Put("/updateUser", UpdateUser)   // Actualizar un usuario
		r.Get("/getUsers", GetUsers)       // Obtener todos los usuarios
		r.Get("/getUser", GetUser)         // Obtener un usuario
		r.Put("/verifyUser", VerifyUser)   // Verificar un usuario
		//TODO manejo de items
		//TODO manejo de prestamos
	})
}

// CreateUser
/*
Recibe los datos de un nuevo usuario, verifica que los campos sean correctos y los guarda en la base de datos
*/
func CreateUser(w http.ResponseWriter, r *http.Request) {

}

// DeleteUser
/*
Recibe el id de un usuario y lo elimina de la base de datos
*/
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// UpdateUser
/*
Recibe los datos de un usuario y los actualiza en la base de datos
*/

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// GetUsers
/*
Obtiene todos los usuarios de la base de datos
*/
func GetUsers(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// GetUser
/*
Obtiene un usuario de la base de datos
*/
func GetUser(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}

// VerifyUser
/*
Recibe el id de un usuario y lo verifica
*/
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	//TODO implementar
}
