package main // import "github.com/xitonix/xhttp/part2"

import (
	"log"
)

func main() {
	server, err := newHTTPServer()
	if err != nil {
		log.Fatal(err)
	}
	server.initRoutes()
	if err = server.start(8080); err != nil {
		log.Fatal(err)
	}
}
