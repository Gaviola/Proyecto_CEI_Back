package models

// Credentials
/*
Estructura de las credenciales de un usuario.
*/
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
