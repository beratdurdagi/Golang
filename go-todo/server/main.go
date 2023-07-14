package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/karalakrepp/Golang/go-react-todo/router"
)

func main() {
	port := os.Getenv("PORT")

	r := router.Router()

	if port == "" {
		port = "8080"
	}

	fmt.Println("Server listening on port: ", port)

	log.Fatal(http.ListenAndServe(":"+port, r))

}
