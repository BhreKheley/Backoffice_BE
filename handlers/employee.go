package handlers

import (
	"absensi-app/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetEmployee retrieves a specific employee by ID
func GetEmployee(c *gin.Context, db *gorm.DB) {
	var employee struct {
		models.Employee
		PositionName string `json:"position_name"`
		DivisionName string `json:"division_name"`
	}
	id := c.Param("id")

	err := db.Table("employee AS e").
		Select("e.*, p.position_name, d.division_name").
		Joins("LEFT JOIN position AS p ON e.position_id = p.id").
		Joins("LEFT JOIN division AS d ON e.division_id = d.id").
		Where("e.id = ?", id).
		Scan(&employee).Error

	if err != nil {
		status := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound || employee.ID == 0 {
			status = http.StatusNotFound
			err = nil
		}
		c.JSON(status, gin.H{
			"status":  "error",
			"message": "Employee not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": employee})
}

// GetAllEmployees retrieves all employees
func GetAllEmployees(c *gin.Context, db *gorm.DB) {
	var employees []struct {
		models.Employee
		PositionName string `json:"position_name"`
		DivisionName string `json:"division_name"`
	}

	err := db.Table("employee AS e").
		Select("e.*, p.position_name, d.division_name").
		Joins("LEFT JOIN position AS p ON e.position_id = p.id").
		Joins("LEFT JOIN division AS d ON e.division_id = d.id").
		Scan(&employees).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve employees"})
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
func CreateEmployee(c *gin.Context, db *gorm.DB) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
		return
	}

	// Check if the associated user is active
	var user models.User
	if err := db.First(&user, "id = ?", employee.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Associated user does not exist",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot create employee for an inactive user",
		})
		return
	}

	// Validate the employee data
	if err := validateEmployee(db, &employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee.CreatedAt = time.Now()

	if err := db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create employee", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": employee})
}

// UpdateEmployee updates an existing employee
func UpdateEmployee(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var employee models.Employee
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		status := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
			err = nil
		}
		c.JSON(status, gin.H{"status": "error", "message": "Employee not found", "error": err.Error()})
		return
	}

	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input data", "error": err.Error()})
		return
	}

	if err := validateEmployee(db, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	employee = input
	employee.UpdatedAt = time.Now()

	if err := db.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update employee", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": employee})
}

// DeleteEmployee deletes an employee by ID
func DeleteEmployee(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var employee models.Employee

	// Retrieve the employee by ID
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Employee not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to retrieve employee",
				"error":   err.Error(),
			})
		}
		return
	}

	// Check if employee is active
	if employee.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Active employees cannot be deleted",
		})
		return
	}

	// Check if the employee has attendance data
	var attendanceCount int64
	if err := db.Model(&models.Attendance{}).Where("user_id = ?", employee.UserID).Count(&attendanceCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to check attendance",
			"error":   err.Error(),
		})
		return
	}

	if attendanceCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Employee has attendance data, cannot be deleted",
		})
		return
	}

	// Check if the associated user is active
	var user models.User
	if err := db.First(&user, "id = ?", employee.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Associated user does not exist, cannot delete employee",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to retrieve associated user",
				"error":   err.Error(),
			})
		}
		return
	}

	if user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Cannot delete employee with an active associated user",
		})
		return
	}

	// Proceed to delete the employee
	if err := db.Delete(&employee).Error; err != nil {
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

// validateEmployee performs validation checks for employee data
func validateEmployee(db *gorm.DB, employee *models.Employee) error {
	// Validate Full Name
	if strings.TrimSpace(employee.Fullname) == "" {
		return errors.New("full name is required and cannot be empty")
	}

	
	// Validate User ID Uniqueness
	if err := db.Where("user_id = ? AND id != ?", employee.UserID, employee.ID).First(&models.Employee{}).Error; err == nil {
		return errors.New("user ID already associated with another employee")
	}
	
	// Validate Full Name Uniqueness
	if err := db.Where("full_name = ? AND id != ?", strings.TrimSpace(employee.Fullname), employee.ID).First(&models.Employee{}).Error; err == nil {
		return errors.New("full name already exists")
	}
	
	// Validate Phone Number
	if strings.TrimSpace(employee.Phone) == "" {
		return errors.New("phone number is required and cannot be empty")
	}
	
	// Validate Division ID
	var division models.Division
	if err := db.First(&division, "id = ?", employee.DivisionID).Error; err != nil {
		return errors.New("invalid division ID")
	}
	
	// Validate Position ID
	var position models.Position
	if err := db.First(&position, "id = ?", employee.PositionID).Error; err != nil {
		return errors.New("invalid position ID")
	}


	// Validate Position-Division Relationship
	if position.DivisionID != employee.DivisionID {
		return errors.New("position does not belong to the specified division")
	}

	return nil
}
