package models

import "database/sql"

// Loan
/*
Estructura de datos de un préstamo.
*/
type Loan struct {
	ID            int
	Status        string
	UserID        int
	AdminID       int
	CreationDate  string
	EndingDate    string
	ReturnDate    string
	Observation   sql.NullString
	Price         float64
	PaymentMethod sql.NullString
}

