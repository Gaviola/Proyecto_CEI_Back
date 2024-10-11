package models

// LoanItem
/*
Estructura de datos de LoanItem, que representa un ítem perteneciente a un préstamo.
*/
type LoanItem struct {
	LoanID int `json:"loan_id"`
	ItemID int `json:"item_id"`
}
