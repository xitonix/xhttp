package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/xitonix/xhttp/part2/storage"
)

type httpServer struct {
	bs *bookStore
	db storage.Database
}

func newHTTPServer() (*httpServer, error) {
	db, err := storage.NewLocalKVStore()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise the database")
	}

	return &httpServer{
		db: db,
		bs: newBookStore(db),
	}, nil
}

func (s *httpServer) initRoutes() {
	http.HandleFunc("/order", s.bs.putOrder(http.MethodPost))
	http.HandleFunc("/list", s.bs.listBooks("all", http.MethodGet))
}

func (s *httpServer) start(port uint32) error {
	if port == 0 {
		port = 80
	}
	defer func() {
		err := s.db.Close()
		if err != nil {
			log.Printf("Failed to close the database connections: %s", err)
		}
	}()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
