package router

import (
	"github.com/gorilla/mux"
	"github.com/karalakrepp/Golang/go-react-todo/middleware"
)

func Router() *mux.Router {

	// router tanımla
	rtr := mux.NewRouter()

	rtr.HandleFunc("/api/task", middleware.GetAllTasks).Methods("GET", "OPTİONS")
	rtr.HandleFunc("/api/task", middleware.CreateTasks).Methods("POST", "OPTİONS")
	rtr.HandleFunc("/api/tasks/{id}", middleware.TaskComplete).Methods("PUT", "OPTIONS")
	rtr.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	rtr.HandleFunc("/api/deleteTask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	rtr.HandleFunc("/api/deleteAllTasks", middleware.DeleteAllTasks).Methods("DELETE", "OPTIONS")

	return rtr
}
