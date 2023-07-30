package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/karalakrepp/Golang/getStudents/db"
	"github.com/karalakrepp/Golang/getStudents/models"
)

type response struct {
	ID  int64  `json:"id,omitempty"`
	Msg string `json:"msg,omitempty"`
}

func GetStudents(w http.ResponseWriter, r *http.Request) {

	students, err := db.GetAllStudents()

	if err != nil {
		log.Fatal("unable to get students %v", err)
	}

	json.NewEncoder(w).Encode(students)
}

func GetStudentByNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	number, err := strconv.Atoi(params["number"])

	if err != nil {
		log.Fatal("Error converting number %v", err)
	}

	student, err := db.GetByNumber(int64(number))

	if err != nil {
		log.Fatal("unable to get student %v", err)
	}

	json.NewEncoder(w).Encode(student)
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {

	var student models.Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Fatal(err)

	}
	student.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.IsAktive = true
	insertNumber := db.InsertStudent(student)

	res := response{
		ID:  insertNumber,
		Msg: "Stock created successfully",
	}
	//verileri json olarak gönder
	json.NewEncoder(w).Encode(res)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	number, err := strconv.Atoi(params["number"])

	if err != nil {
		log.Fatal(err)
	}
	var student models.Student

	err = json.NewDecoder(r.Body).Decode(&student)

	if err != nil {
		log.Fatal("unable to decode student")
		return
	}

	student.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	rowsAffected := db.UpdateByNumber(int64(number), student)

	resp := response{
		ID:  rowsAffected,
		Msg: "Updated student successfully",
	}
	json.NewEncoder(w).Encode(&resp)

}

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	number, err := strconv.Atoi(params["number"])

	if err != nil {
		log.Fatal(err)
	}
	var student models.Student
	student.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	rowsAffected, err := db.DeactivateStudent(int64(number))

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	resp := response{
		ID:  rowsAffected,
		Msg: "Deleted student successfully",
	}
	json.NewEncoder(w).Encode(&resp)

}
