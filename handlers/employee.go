package handlers

import (
	"absensi-app/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetEmployee retrieves an employee by ID
func GetEmployee(c *gin.Context, db *sqlx.DB) {
	var employee models.Employee
	id := c.Param("id")

	err := db.QueryRow("SELECT id, user_id, full_name, phone, position_id, division_id, created_at, updated_at FROM employee WHERE id = $1", id).
		Scan(&employee.ID, &employee.UserID, &employee.Fullname, &employee.Phone, &employee.PositionID, &employee.DivisionID, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Employee not found",
		})
		return
	}

	// Format response to match the user response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":          employee.ID,
			"user_id":     employee.UserID,
			"full_name":   employee.Fullname,
			"phone":       employee.Phone,
			"position_id": employee.PositionID,
			"division_id": employee.DivisionID,
			"created_at":  employee.CreatedAt,
			"updated_at":  employee.UpdatedAt,
		},
	})
}

// GetAllEmployees retrieves all employees
func GetAllEmployees(c *gin.Context, db *sqlx.DB) {
	var employees []models.Employee

	err := db.Select(&employees, "SELECT id, user_id, full_name, phone, position_id, division_id, is_active, created_at, updated_at FROM employee")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve employees",
			"error":   err.Error(), // Tambahkan detail error
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

	// Format response to match the user response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   formatEmployees(employees),
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

	// Check if the user_id is already associated with an employee
	var existingEmployeeCount int
	err := db.Get(&existingEmployeeCount, "SELECT COUNT(*) FROM employee WHERE user_id = $1", employee.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to check existing employee",
			"error":   err.Error(),
		})
		return
	}

	if existingEmployeeCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ini sudah terkait dengan employee lain",
		})
		return
	}

	// Set timestamps
	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	// Insert new employee record
	_, err = db.Exec("INSERT INTO employee (user_id, full_name, phone, position_id, division_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		employee.UserID, employee.Fullname, employee.Phone, employee.PositionID, employee.DivisionID, employee.CreatedAt, employee.UpdatedAt)
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
		"data": gin.H{
			"user_id":     employee.UserID,
			"full_name":   employee.Fullname,
			"phone":       employee.Phone,
			"position_id": employee.PositionID,
			"division_id": employee.DivisionID,
		},
	})
}

// UpdateEmployee updates an existing employee by ID
func UpdateEmployee(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var employee models.Employee

	// Check if employee exists
	err := db.Get(&employee, "SELECT * FROM employee WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Employee not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	employee.UpdatedAt = time.Now()

	_, err = db.Exec("UPDATE employee SET full_name = $1, phone = $2, position_id = $3, division_id = $4, updated_at = $5 WHERE id = $6",
		employee.Fullname, employee.Phone, employee.PositionID, employee.DivisionID, employee.UpdatedAt, id)

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

	// Check if employee exists
	var employee models.Employee
	err := db.Get(&employee, "SELECT * FROM employee WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Employee tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal mendapatkan data employee",
				"error":   err.Error(),
			})
		}
		return
	}

	// Check if employee has related attendance data
	var attendanceCount int
	err = db.Get(&attendanceCount, "SELECT COUNT(*) FROM attendance WHERE user_id = $1", employee.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memeriksa data absensi",
			"error":   err.Error(),
		})
		return
	}

	// If employee has attendance data, prevent deletion
	if attendanceCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Employee memiliki data absensi, tidak dapat dihapus",
		})
		return
	}

	// Check if the related user is inactive
	var userIsActive bool
	err = db.Get(&userIsActive, "SELECT is_active FROM \"user\" WHERE id = $1", employee.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memeriksa status user",
			"error":   err.Error(),
		})
		return
	}

	// Only delete employee if related user is inactive
	if userIsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User terkait masih aktif, employee tidak dapat dihapus",
		})
		return
	}

	// Delete employee if no related attendance and user is inactive
	_, err = db.Exec("DELETE FROM employee WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus employee",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Employee berhasil dihapus",
	})
}

// Helper function to format employees
func formatEmployees(employees []models.Employee) []gin.H {
	var formatted []gin.H
	for _, emp := range employees {
		formatted = append(formatted, gin.H{
			"id":          emp.ID,
			"user_id":     emp.UserID,
			"full_name":   emp.Fullname,
			"phone":       emp.Phone,
			"position_id": emp.PositionID,
			"division_id": emp.DivisionID,
			"is_active":   emp.IsActive,
			"created_at":  emp.CreatedAt,
			"updated_at":  emp.UpdatedAt,
		})
	}
	return formatted
}
