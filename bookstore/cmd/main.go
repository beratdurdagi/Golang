package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karalakrepp/Golang/bookstore/pkg/routes"
	_ "gorm.io/driver/postgres"
)

func main() {
	router := mux.NewRouter()
	routes.RegisterBookStoreRoutes(router)
	http.Handle("/", router)
	fmt.Println("server connected")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
