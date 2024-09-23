package models

type Loan struct {
	ID            int
	status        string
	userID        int
	adminID       int
	creationDate  string
	endingDate    string
	returnDate    string
	observation   string
	price         float64
	paymentMethod string
}
