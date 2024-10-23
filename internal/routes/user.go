package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	// "github.com/Gaviola/Proyecto_CEI_Back.git/internal/middlewares"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/middlewares"
	"github.com/Gaviola/Proyecto_CEI_Back.git/internal/repositories"
	"github.com/Gaviola/Proyecto_CEI_Back.git/models"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware) // Middleware de verificación de token

		// Rutas para usuarios
		r.Route("/{userID}", func(r chi.Router) {
			r.Post("/createLoan", createLoan)   // Crear un préstamo
			r.Delete("/cancelLoan", cancelLoan) // Cancelar prestamo
			r.Get("/getLoans", getLoans)        // Obtener préstamos
			r.Get("/getItems", getItems)        // Obtener ítems
			r.Patch("/updateUser", updateUser)  // Actualizar datos del usuario
			r.Get("/getUser", getUser)          // Obtener datos del usuario

		})
	})
}

// CreaLoan
/*
Crea un préstamo de un ítem para el usuario.
*/
func createLoan(w http.ResponseWriter, r *http.Request) {
	idUser := chi.URLParam(r, "userID")

	loan := models.Loan{}
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	loan.UserID, err = strconv.Atoi(idUser)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = repositories.DBSaveLoan(loan)
	if err != nil {
		http.Error(w, "Error creating loan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// CancelLoan
/*
Cancela un préstamo de un ítem para el usuario.
*/
func cancelLoan(w http.ResponseWriter, r *http.Request) {
	idUser := chi.URLParam(r, "userID")

	loan := models.Loan{}
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	loan.UserID, err = strconv.Atoi(idUser)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = repositories.DBDeleteLoan(loan.ID)
	if err != nil {
		http.Error(w, "Error deleting loan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetLoans
/*
Obtiene todos los préstamos del usuario.
*/
func getLoans(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	loans, err := repositories.DBGetLoansByUserID(id)
	if err != nil {
		http.Error(w, "Error getting loans", http.StatusInternalServerError)
		return
	}
	// Devolver los préstamos en formato JSON
	err = json.NewEncoder(w).Encode(loans)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

// GetItems
/*
Obtiene todos los ítems disponibles para prestar a un usuario.
*/
func getItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item

	items, err := repositories.DBGetAvailableItems()
	if err != nil {
		http.Error(w, "Error getting items", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

// UpdateUser
/*
Actualiza los datos del usuario. Datos como ID o email no se pueden modificar.
*/
func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	updatedUser := models.User{}
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	//Chequeo que no se intente modificar ID o email
	if updatedUser.ID != id || updatedUser.Email != "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = repositories.DBUpdateUser(id, updatedUser)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// GetUser
/*
Obtiene los datos del usuario.
*/
func getUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := repositories.DBGetUserByID(id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)

}
