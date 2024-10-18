package routes

import (
	"absensi-app/handlers"
	"absensi-app/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SetupRoutes initializes all application routes
func SetupRoutes(r *gin.Engine, db *sqlx.DB) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth Routes
	auth := r.Group("")
	{
		auth.POST("/auth", func(c *gin.Context) {
			handlers.Login(c, db)
		})
	}

	// Routes with auth middleware
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Attendance Routes
	attendance := protected.Group("/attendance")
	{
		attendance.POST("/checkin", func(c *gin.Context) {
			handlers.CheckIn(c, db)
		})
		attendance.POST("/checkout", func(c *gin.Context) {
			handlers.CheckOut(c, db)
		})
	}

	// Employee Routes
	employee := protected.Group("/employee")
	{
		employee.GET("/:id", func(c *gin.Context) {
			handlers.GetEmployee(c, db)
		})
		employee.POST("/", func(c *gin.Context) {
			handlers.CreateEmployee(c, db)
		})
	}

	// User Routes
	user := protected.Group("/user")
	{
		user.GET("/:id", func(c *gin.Context) {
			handlers.GetUserByID(c, db)
		})
		user.GET("/", func(c *gin.Context) {
			handlers.GetAllUsers(c, db)
		})
	}
}
