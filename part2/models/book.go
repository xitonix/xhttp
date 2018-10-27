package models

import "time"

// Book encapsulates the details of a book
type Book struct {
	// Title the book title
	Title string `json:"title"`
	// Category the category to which the book belongs
	Category string `json:"category"`
	// Author the author's name
	Author string `json:"author"`
	// Published publish date
	Published time.Time `json:"published"`
	// Pages number of pages
	Pages int `json:"pages"`
	// Catalogued the date when the book was added to Mr. Bookworm's catalogue
	Catalogued time.Time `json:"catalogued"`
	// Stock number of the copies available in stock
	Stock int `json:"stock"`
	// ISBN International Standard Book Number
	ISBN string `json:"-"`
}
