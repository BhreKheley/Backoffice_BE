package handlers

import (
	"absensi-app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetEmployee retrieves an employee by ID
func GetEmployee(c *gin.Context, db *sqlx.DB) {
	var employee models.Employee
	id := c.Param("id")

	err := db.QueryRow("SELECT id, user_id, full_name, phone, position, created_at, updated_at FROM employee WHERE id = $1", id).
		Scan(&employee.ID, &employee.UserID, &employee.Fullname, &employee.Phone, &employee.Position, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// CreateEmployee creates a new employee
func CreateEmployee(c *gin.Context, db *sqlx.DB) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	_, err := db.Exec("INSERT INTO employee (user_id, full_name, phone, position, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		employee.UserID, employee.Fullname, employee.Phone, employee.Position, employee.CreatedAt, employee.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee created successfully", "data": employee})
}
