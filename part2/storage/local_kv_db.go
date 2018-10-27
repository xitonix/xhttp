package storage

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/xitonix/xhttp/part2/models"
)

type LocalKeyValueDB struct {
	books  *badger.DB
	orders *badger.DB
}

// NewLocalKVStore creates a new instance of a local key-value store database
func NewLocalKVStore() (*LocalKeyValueDB, error) {
	books, err := openDatabase("books")
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise the Books database")
	}
	orders, err := openDatabase("orders")
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise the Orders database")
	}
	return &LocalKeyValueDB{
		books:  books,
		orders: orders,
	}, nil
}

// ReadCategory returns all the books in the specified category.
func (l *LocalKeyValueDB) ReadCategory(category string) ([]models.Book, error) {
	books := make([]models.Book, 0)
	err := l.books.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			isbn := item.Key()
			b, err := item.Value()
			if err != nil {
				return errors.Wrap(err, "failed to read the book data from the database")
			}
			book := models.Book{}
			err = json.Unmarshal(b, &book)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal the book details")
			}
			if strings.EqualFold(book.Category, category) || strings.EqualFold(category, "all") {
				book.ISBN = string(isbn)
				books = append(books, book)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the new arrivals")
	}
	return books, nil
}

// PlaceOrder places a new order and returns order ID if succeeded.
func (l *LocalKeyValueDB) PlaceOrder(order models.Order) (string, error) {
	id := uuid.New()
	order.Date = time.Now().UTC()
	err := l.orders.Update(func(txn *badger.Txn) error {
		b, err := json.Marshal(order)
		if err != nil {
			return errors.Wrap(err, "failed to marshal the order into Json")
		}
		return txn.Set([]byte(id.String()), b)
	})
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// Close closes the local database file
func (l *LocalKeyValueDB) Close() error {
	err := l.books.Close()
	if err != nil {
		err = errors.Wrap(err, "failed to close the books database")
	}
	e := l.orders.Close()
	if e != nil {
		err = errors.Errorf("failed to close the orders database: %s. %s", e, err)
	}
	return err
}

func openDatabase(path string) (*badger.DB, error) {
	opts := badger.DefaultOptions
	path = filepath.Join(".", "data", path)
	opts.Dir = path
	opts.ValueDir = path
	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open the database")
	}
	return db, nil
}
