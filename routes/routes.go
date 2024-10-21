package routes

import (
	"absensi-app/handlers"
	"absensi-app/middleware"
	"net/http"

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
		// Create User Route (now inside the protected group)
		protected.POST("/create_user", func(c *gin.Context) {
			handlers.CreateUser(c, db)
		})
	}

	// Add auth route with middleware
	authProtected := protected.Group("/auth")
	{
		authProtected.GET("/get-user-by-token", func(c *gin.Context) {
			handlers.GetUserByToken(c, db)
		})
	}

	// Route untuk menampilkan daftar semua rute
	r.GET("/list-routes", func(c *gin.Context) {
		routes := r.Routes()
		var routeList []map[string]string

		// Loop melalui daftar rute yang ada dan tambahkan ke array
		for _, route := range routes {
			routeInfo := map[string]string{
				"method": route.Method,
				"path":   route.Path,
			}
			routeList = append(routeList, routeInfo)
		}

		// Kembalikan daftar rute dalam bentuk JSON
		c.JSON(http.StatusOK, gin.H{
			"routes": routeList,
		})
	})
}
