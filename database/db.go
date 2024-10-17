package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// InitDB initializes the database connection
func InitDB() *sqlx.DB {
	// Update your connection string accordingly
	connStr := "host=postgres user=postgres dbname=absensi_app password=mysecretpassword sslmode=disable"

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	DB = db
	fmt.Println("Database connection established")
	return DB
}
