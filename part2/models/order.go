package models

import "time"

// Order represents an order to purchase a book
type Order struct {
	//ID the order ID
	ID string `json:"-"`
	// Count number of copies to order
	Count int `json:"count"`
	// ISBN the ISBN of the book to order
	ISBN string `json:"isbn"`
	// Date the date when the order has been received
	Date time.Time `json:"date"`
}
