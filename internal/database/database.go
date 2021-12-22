package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewDatabse : return pointer to db object
func NewDatabase() (*gorm.DB, error) {
	fmt.Println("Setting up DB connection")

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=America/New_York", dbUsername, dbPassword, dbHost, dbPort, dbTable)

	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		return db, err
	}

	if err := db.DB().Ping(); err != nil {
		return db, err
	}

	return db, nil
}
