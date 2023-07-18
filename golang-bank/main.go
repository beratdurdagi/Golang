package main

import (
	"log"

	"github.com/karalakrepp/Golang/golang-bank/API"
	"github.com/karalakrepp/Golang/golang-bank/models"
)

func main() {

	store, err := models.NewPostgresStrore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := API.NewAPIServer(":8000", store)
	server.Run()
}
