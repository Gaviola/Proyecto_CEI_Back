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
