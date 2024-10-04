package models

// Item
/*
Estructura de datos de un ítem.
*/
type Item struct {
	ID       int
	ItemType string
	ItemTypeID int
	Code     string
	Price    float64
}
