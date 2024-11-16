package main

import (
	"absensi-app/database"
	"absensi-app/models"
	"absensi-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi koneksi database
	db := database.InitDB()

	// Auto-migrate untuk membuat atau memperbarui tabel berdasarkan model
	err := db.AutoMigrate(
		&models.Attendance{},
		&models.Division{},
		&models.Employee{},
		&models.Permission{},
		&models.Position{},
		&models.Role{},
		&models.RolePermission{},
		&models.Status{},
		&models.User{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, db)

	// Menjalankan server di port 8080
	r.Run(":8080")
}
