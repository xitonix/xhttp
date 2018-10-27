package storage

import (
	"io"

	"github.com/xitonix/xhttp/part2/models"
)

// Database is the interface for our bookstore storage
type Database interface {
	ReadCategory(category string) ([]models.Book, error)
	PlaceOrder(order models.Order) (string, error)
	io.Closer
}
