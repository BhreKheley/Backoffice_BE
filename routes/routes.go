package routes

import (
	"absensi-app/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SetupRoutes initializes all application routes
func SetupRoutes(r *gin.Engine, db *sqlx.DB) {
	// Auth Routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			handlers.Login(c, db)
		})
	}

	// Attendance Routes
	attendance := r.Group("/attendance")
	{
		attendance.POST("/checkin", func(c *gin.Context) {
			handlers.CheckIn(c, db)
		})
		attendance.POST("/checkout", func(c *gin.Context) {
			handlers.CheckOut(c, db)
		})
	}

	// Employee Routes
	employee := r.Group("/employee")
	{
		employee.GET("/:id", func(c *gin.Context) {
			handlers.GetEmployee(c, db)
		})
		employee.POST("/", func(c *gin.Context) {
			handlers.CreateEmployee(c, db)
		})
	}

	// User Routes
	user := r.Group("/user")
	{
		user.POST("/create", func(c *gin.Context) {
			handlers.CreateUser(c, db)
		})
		user.GET("/:id", func(c *gin.Context) {
			handlers.GetUserByID(c, db)
		})
		user.GET("/", func(c *gin.Context) {
			handlers.GetAllUsers(c, db)
		})
	}
}
