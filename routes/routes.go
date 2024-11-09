package routes

import (
	"absensi-app/handlers"
	"absensi-app/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes initializes all application routes using GORM
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/check-username", func(c *gin.Context) { handlers.CheckUsername(c, db) })
	r.POST("/check-email", func(c *gin.Context) { handlers.CheckEmail(c, db) })

	// Auth Routes
	auth := r.Group("auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			handlers.Login(c, db)
		})
	}

	// Routes with auth middleware
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Attendance Routes
	attendance := protected.Group("/attendance")
	{
		attendance.Use(middleware.CheckPermission("VIEW_ATTENDANCE", db))
		attendance.GET("/", func(c *gin.Context) {
			handlers.GetAllAttendance(c, db)
		})
		attendance.GET("/:id", func(c *gin.Context) {
			handlers.GetAttendanceByID(c, db)
		})
		attendance.Use(middleware.CheckPermission("MANAGE_ATTENDANCE", db))
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
		employee.Use(middleware.CheckPermission("VIEW_EMPLOYEES", db))
		employee.GET("/:id", func(c *gin.Context) {
			handlers.GetEmployee(c, db)
		})
		employee.GET("/", func(c *gin.Context) {
			handlers.GetAllEmployees(c, db)
		})
		employee.Use(middleware.CheckPermission("MANAGE_EMPLOYEES", db))
		employee.POST("/create_employee", func(c *gin.Context) {
			handlers.CreateEmployee(c, db)
		})
		employee.PUT("/:id", func(c *gin.Context) {
			handlers.UpdateEmployee(c, db)
		})
		employee.DELETE("/:id", func(c *gin.Context) {
			handlers.DeleteEmployee(c, db)
		})
	}

	// User Routes
	user := protected.Group("/user")
	{
		user.Use(middleware.CheckPermission("VIEW_USERS", db))
		user.GET("/:id", func(c *gin.Context) {
			handlers.GetUserByID(c, db)
		})
		user.GET("/", func(c *gin.Context) {
			handlers.GetAllUsers(c, db)
		})
		user.Use(middleware.CheckPermission("MANAGE_USERS", db))
		user.POST("/create_user", func(c *gin.Context) {
			handlers.CreateUser(c, db)
		})
		user.PUT("/:id", func(c *gin.Context) {
			handlers.UpdateUser(c, db)
		})
		user.DELETE("/:id", func(c *gin.Context) {
			handlers.DeleteUser(c, db)
		})
	}

	// Role Routes
	role := protected.Group("/role")
	{
		role.Use(middleware.CheckPermission("VIEW_USERS", db))
		role.GET("/", func(c *gin.Context) {
			handlers.GetRoles(c, db)
		})
		role.POST("/create_role", func(c *gin.Context) {
			handlers.CreateRole(c, db)
		})
	}

	// Permission Routes
	permission := protected.Group("/permission")
	{
		permission.GET("/", func(c *gin.Context) {
			handlers.GetPermissions(c, db)
		})
		permission.GET("/byrole", func(c *gin.Context) {
			handlers.GetPermissionsByRole(c, db)
		})
		permission.POST("/create_permission", func(c *gin.Context) {
			handlers.CreatePermission(c, db)
		})
		permission.POST("/assign_permission_to_role", func(c *gin.Context) {
			handlers.AssignPermissionToRole(c, db)
		})
	}

	// Division Routes
	division := protected.Group("/division")
	{
		division.Use(middleware.CheckPermission("VIEW_DIVISIONS", db))
		division.GET("/", func(c *gin.Context) {
			handlers.GetAllDivisions(c, db)
		})
		division.GET("/:id", func(c *gin.Context) {
			handlers.GetDivisionByID(c, db)
		})
		division.Use(middleware.CheckPermission("MANAGE_DIVISIONS", db))
		division.POST("/create", func(c *gin.Context) {
			handlers.CreateDivision(c, db)
		})
		division.PUT("/:id", func(c *gin.Context) {
			handlers.UpdateDivision(c, db)
		})
		division.DELETE("/:id", func(c *gin.Context) {
			handlers.DeleteDivision(c, db)
		})
	}

	// Position Routes
	position := protected.Group("/position")
	{
		position.Use(middleware.CheckPermission("VIEW_POSITIONS", db))
		position.GET("/", func(c *gin.Context) {
			handlers.GetAllPositions(c, db)
		})
		position.GET("/:id", func(c *gin.Context) {
			handlers.GetPositionByID(c, db)
		})
		position.Use(middleware.CheckPermission("MANAGE_POSITIONS", db))
		position.POST("/create", func(c *gin.Context) {
			handlers.CreatePosition(c, db)
		})
		position.PUT("/:id", func(c *gin.Context) {
			handlers.UpdatePosition(c, db)
		})
		position.DELETE("/:id", func(c *gin.Context) {
			handlers.DeletePosition(c, db)
		})
	}

	// Auth route with middleware for getting user by token
	authProtected := protected.Group("/auth")
	{
		authProtected.GET("/get-user-by-token", func(c *gin.Context) {
			handlers.GetUserByToken(c, db)
		})
	}

	// Route to display all registered routes
	r.GET("/list-routes", func(c *gin.Context) {
		routes := r.Routes()
		var routeList []map[string]string

		for _, route := range routes {
			routeInfo := map[string]string{
				"method": route.Method,
				"path":   route.Path,
			}
			routeList = append(routeList, routeInfo)
		}

		c.JSON(http.StatusOK, gin.H{
			"routes": routeList,
		})
	})
}
