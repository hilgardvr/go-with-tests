package main

import (
	"log"
	"net/http"
	"github.com/hilgardvr/go-with-tests/application/application"
)

func main() {
	server := application.NewPlayerServer(application.NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}