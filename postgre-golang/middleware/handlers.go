package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/karalakrepp/Golang/postgres/models"
	_ "github.com/lib/pq"
)

type response struct {
	ID  int64  `json:"id,omitempty"`
	Msg string `json:"msg,omitempty"`
}

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

func NewStock(w http.ResponseWriter, r *http.Request) {
	// stock oluştur
	var stock models.Stock

	// gelen jsonları decode et ve err handling yap
	err := json.NewDecoder(r.Body).Decode(&stock)

	//decode hatası varsa yazdır
	if err != nil {
		log.Fatal("unable to decode the request body %v", err)

	}
	//database tarafındaki id yi eşle
	insertID := insertStock(stock)

	//verilecek mesajı belirle
	res := response{
		ID:  insertID,
		Msg: "Stock created successfully",
	}
	//verileri json olarak gönder
	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {

	// parametreleri maple ve params eşitle
	params := mux.Vars(r)

	//id ve err tanımla
	id, err := strconv.Atoi(params["id"])

	//err handling yap
	if err != nil {
		log.Fatal("unable to conver string into int %v", err)
	}

	//databasedens stock bilgisini eşle
	stock, err := getStock(int64(id))

	// error handling yap
	if err != nil {
		log.Fatalf("unable to get stock %v", err)
	}

	// response gönder
	json.NewEncoder(w).Encode(stock)

}

func GetStocks(w http.ResponseWriter, q *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatal("Unable to gel all stocks %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UptadeStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatal("unable to conver string into int %v", err)

	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("unable to decode the request body %v", err)
	}

	updateRows := uptadeStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated Successfully %v", updateRows)

	res := response{
		ID:  int64(id),
		Msg: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatal("unable to conver string into int %v", err)

	}

	deletedRows := deletedStock(int64(id))

	msg := fmt.Sprintf("Stock Removed Successfully %v", deletedRows)

	res := response{
		ID:  int64(id),
		Msg: msg,
	}

	json.NewEncoder(w).Encode(res)

}

//********************************************************************************************//
// handler function for database

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name,price,company) VALUES($1, $2, $3) RETURNING stockid`

	var id int64
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatal("unable to execute the query %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

func getStock(stockId int64) (models.Stock, error) {
	db := CreateConnection()

	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid =$1`

	row := db.QueryRow(sqlStatement, stockId)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	// error handling
	case sql.ErrNoRows:
		log.Fatal("no rows were returned!  %v", err)
		return stock, nil
	case nil:
		return stock, nil

	default:
		log.Fatal("Unable to scan the row. %v", err)
	}

	return stock, err

}

func getAllStocks() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlStatement := "SELECT * FROM stocks "

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatal("unable to query %v", err)

	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatal("Unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)

	}
	return stocks, err

}

func uptadeStock(stockId int64, stock models.Stock) int64 {

	db := CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE  stocks SET name=$2,price=$3,company=$4 WHERE stockid=$1`

	result, err := db.Exec(sqlStatement, stockId, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatal("unable to execute the query %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatal("Error while getting rows affected %v", err)
	}

	return rowsAffected

}

func deletedStock(stockId int64) int64 {

	// create the postgres db connection
	db := CreateConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, stockId)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
