package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection using GORM
func InitDB() *gorm.DB {
	// Update your connection string accordingly
	connStr := "host=postgres user=postgres dbname=absensi_app password=mysecretpassword sslmode=disable"

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	DB = db
	fmt.Println("Database connection established using GORM")
	return DB
}
