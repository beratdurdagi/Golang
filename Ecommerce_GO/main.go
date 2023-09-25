package main

import (
	"log"

	"github.com/Karalakrepp/Golang/Ecommerce_GO/routes"
)

func main() {

	server := routes.NewAPIServer(log.Default(), ":8080")

	server.Run()

}
