package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xitonix/xhttp/part2/models"
	"github.com/xitonix/xhttp/part2/storage"
)

const categoryKey = "category"

type bookStore struct {
	db storage.Database
}

func newBookStore(db storage.Database) *bookStore {
	return &bookStore{
		db: db,
	}
}

func (b *bookStore) putOrder(allowedMethod string) http.HandlerFunc {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != allowedMethod {
			http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(w, "invalid order data", err, http.StatusBadRequest)
			return
		}
		order := models.Order{}
		err = json.Unmarshal(body, &order)
		if err != nil {
			writeError(w, "failed to read the order", err, http.StatusInternalServerError)
			return
		}
		id, err := b.db.PlaceOrder(order)
		if err != nil {
			writeError(w, "failed to place a new order", err, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, id)
	}
	return logger(handler)
}

func (b *bookStore) listBooks(defaultCategory, allowedMethod string) http.HandlerFunc {
	handler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method != allowedMethod {
			http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
			return
		}
		cat := req.URL.Query().Get(categoryKey)
		if len(cat) == 0 {
			cat = defaultCategory
		}
		books, err := b.db.ReadCategory(cat)
		if err != nil {
			writeError(w, "failed to fetch the books", err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err = enc.Encode(books)
		if err != nil {
			writeError(w, "failed to serialise the book list", err, http.StatusInternalServerError)
		}
	}
	return logger(handler)
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Servinge %s %s %s", r.Method, r.URL.Path, r.URL.Query())
		next(w, r)
	}
}

func writeError(w http.ResponseWriter, msg string, err error, code int) {
	log.Printf("%s: %s", msg, err)
	http.Error(w, msg, code)
}
