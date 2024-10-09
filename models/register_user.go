package models

import "reflect"

type RegisterUser struct {
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	StudentId  int    `json:"student_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      int    `json:"phone"`
	Role       string `json:"role"`
	Dni        int    `json:"dni"`
	CreatorId  int    `json:"creator_id"`
	School     string `json:"school"`
	IsVerified bool   `json:"is_verified"`
	Hash       []byte
}

// IsEmpty
/*
Devuelve true si el usuario es vacio, es decir, si no tiene ningun campo con informacion.
Devuelve false si el usuario tiene al menos un campo con informacion.
*/
func (user *RegisterUser) IsEmpty() bool {
	value := reflect.ValueOf(user).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.IsZero() {
			return false
		}
	}
	return true
}
