package models

// Item
/*
Estructura de datos de un ítem.
*/
type Item struct {
	ID         int 		`json:"id"`
	ItemType   string	`json:"item_type"`
	ItemTypeID int     `json:"item_type_id"`
	Code       string  `json:"code"`
	Price      float64 `json:"price"`
}
