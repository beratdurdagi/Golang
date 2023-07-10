package router

import (
	"github.com/gorilla/mux"
	"github.com/karalakrepp/Golang/postgres/middleware"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock", middleware.GetStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newstock", middleware.NewStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stock/{id}", middleware.UptadeStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/removestock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

	return router
}
