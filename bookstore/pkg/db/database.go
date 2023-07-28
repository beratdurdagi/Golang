package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	dsn := "host=localhost user=postgres password=12345678 dbname=bookstore port=5433 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	db = database
}
func GetDB() *gorm.DB {
	return db
}
