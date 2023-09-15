package main

import (
	"log"

	database "github.com/karalakrepp/Golang/JWT/Database"
)

func main() {
	store, err := database.NewPostgresStrore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewServer(":8000", store)
	
	server.Run()
}
