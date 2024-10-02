package models

// ItemType
/*
Estructura de datos de un tipo de ítem.
*/
type ItemType struct {
	ID        int   `json:"id"`
	Name      string `json:"name"`
	IsGeneric bool  `json:"is_generic"`
}
