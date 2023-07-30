package router

import (
	"github.com/gorilla/mux"
	"github.com/karalakrepp/Golang/getStudents/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/student/{number}", controller.GetStudentByNumber).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/studentbyname/{first_name}", controller.GetStudentByNumber).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/student", controller.GetStudents).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newstudent", controller.CreateStudent).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/student/{number}", controller.UpdateStudent).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/removestudent/{number}", controller.DeleteStudent).Methods("DELETE", "OPTIONS")

	return router
}
