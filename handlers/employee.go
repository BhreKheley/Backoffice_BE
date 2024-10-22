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
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Employee not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   employee,
	})
}

// GetAllEmployees retrieves all employees
func GetAllEmployees(c *gin.Context, db *sqlx.DB) {
	var employees []models.Employee

	err := db.Select(&employees, "SELECT id, user_id, full_name, phone, position, created_at, updated_at FROM employee")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve employees",
		})
		return
	}

	if len(employees) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No employees found",
			"data":    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   employees,
	})
}

// CreateEmployee creates a new employee
func CreateEmployee(c *gin.Context, db *sqlx.DB) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	_, err := db.Exec("INSERT INTO employee (user_id, full_name, phone, position, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		employee.UserID, employee.Fullname, employee.Phone, employee.Position, employee.CreatedAt, employee.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create employee",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Employee created successfully",
		"data":    employee,
	})
}

// UpdateEmployee updates an existing employee by ID
func UpdateEmployee(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	employee.UpdatedAt = time.Now()

	_, err := db.Exec("UPDATE employee SET full_name = $1, phone = $2, position = $3, updated_at = $4 WHERE id = $5",
		employee.Fullname, employee.Phone, employee.Position, employee.UpdatedAt, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update employee",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Employee updated successfully",
	})
}

// DeleteEmployee deletes an employee by ID
func DeleteEmployee(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM employee WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete employee",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Employee deleted successfully",
	})
}
