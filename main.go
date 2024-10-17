package main

import (
	"absensi-app/database"
	"absensi-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.InitDB()
	r := gin.Default()

	// Setup all routes
	routes.SetupRoutes(r, db)

	r.Run(":8080") // Menjalankan server di port 8080
}
