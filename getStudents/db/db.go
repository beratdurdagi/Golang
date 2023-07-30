package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/karalakrepp/Golang/getStudents/models"
	_ "github.com/lib/pq"
)

func CreateConnection() *sql.DB {

	// .env dosyasına bağlantı gönderiyoruz
	err := godotenv.Load(".env")

	// .env dosyası için error handling
	if err != nil {
		log.Fatal("Error loading .env file")

	}

	//sql bağlanıyoruz

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	// bağlı olup olmadığını kontrol için error handling
	err = db.Ping()

	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	// hata yoksa bu mesaj dönecek
	fmt.Println("Successfully created connection to postgres")

	//return ediyoruz

	return db

}

func InsertStudent(student models.Student) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := " INSERT INTO student(first_name,last_name,phone,number,email,address,created_at,updated_at,isaktive) VALUES($1, $2, $3,$4,$5,$6,$7,$8,$9) RETURNING number"

	var number int64

	err := db.QueryRow(sqlStatement,

		student.FirstName,
		student.LastName,
		student.Phone,
		student.Number,
		student.Email,
		student.Address,
		student.CreatedAt,
		student.UpdatedAt,
		student.IsAktive).Scan(&number)

	if err != nil {
		log.Fatal(err)

	}

	fmt.Printf("Inserted a single record %v", number)
	return number

}

func GetAllStudents() ([]models.Student, error) {

	db := CreateConnection()
	defer db.Close()

	var students []models.Student
	query := "Select * from student"

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal("unable to query %v", err)

	}

	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Phone,
			&student.Number,
			&student.Email,
			&student.Address,
			&student.CreatedAt,
			&student.UpdatedAt,
			&student.IsAktive)

		if err != nil {
			log.Fatal("unable to scan row %v", err)
		}
		students = append(students, student)
	}
	return students, err

}

func GetByNumber(number int64) (models.Student, error) {
	db := CreateConnection()
	defer db.Close()

	var student models.Student
	query := "SELECT * FROM student WHERE number=$1"

	row := db.QueryRow(query, number)

	err := row.Scan(&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Phone,
		&student.Number,
		&student.Email,
		&student.Address,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.IsAktive)

	switch err {
	// error handling
	case sql.ErrNoRows:
		log.Fatal("no rows were returned!  %v", err)
		return student, nil
	case nil:
		return student, nil

	default:
		log.Fatal("Unable to scan the row. %v", err)
	}

	return student, err

}

func UpdateByNumber(number int64, student models.Student) int64 {
	db := CreateConnection()
	defer db.Close()

	query := "UPDATE student SET first_name=$2,last_name=$3,phone=$4,email=$5,address=$6 WHERE number=$1"

	result, err := db.Exec(query, number, student.FirstName, student.LastName, student.Phone, student.Email, student.Address)

	if err != nil {
		log.Fatal("unable to execute the query %v", err)
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal("Error while getting rows affected %v", err)
	}

	return rowsAffected

}
func DeactivateStudent(number int64) (int64, error) {
	db := CreateConnection()
	defer db.Close()

	query := "UPDATE student SET isaktive=false WHERE number=$1"

	result, err := db.Exec(query, number)

	if err != nil {
		log.Fatal("unable to execute the query %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal("Error while getting rows affected %v", err)
	}

	return rowsAffected, nil
}
