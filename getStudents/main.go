package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	router "github.com/karalakrepp/Golang/getStudents/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	rtr := router.Router()

	fmt.Println("Starting server on port ", port)

	log.Fatal(http.ListenAndServe(":"+port, rtr))
}
