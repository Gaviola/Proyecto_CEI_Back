package models

import "reflect"

// User
/*
Estructura de datos que representa a un usuario.
*/
type User struct {
	ID         int
	Name       string
	Lastname   string
	StudentId  int //legajo
	Email      string
	Phone      int
	Role       string
	Dni        int
	CreatorId  int
	School     string
	IsVerified bool
	Hash       []byte
}

// IsEmpty
/*
Devuelve true si el usuario es vacio, es decir, si no tiene ningun campo con informacion.
Devuelve false si el usuario tiene al menos un campo con informacion.
*/
func (user *User) IsEmpty() bool {
	value := reflect.ValueOf(user).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.IsZero() {
			return false
		}
	}
	return true
}
