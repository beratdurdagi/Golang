package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	controller "github.com/karalakrepp/golang/controllers"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()

	router.HandleFunc("/movies", controller.GetMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", controller.GetMovie).Methods("GET")
	router.HandleFunc("/movies", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", controller.UpdateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", controller.DeleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port  %s\n ", port)

	log.Fatal(http.ListenAndServe(":"+port, router))

}
