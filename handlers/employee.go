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

// Helper function to handle errors uniformly
func handleError(c *gin.Context, status int, message string, err error) {
	if err != nil {
		c.JSON(status, gin.H{
			"status":  "error",
			"message": message,
			"error":   err.Error(),
		})
	} else {
		c.JSON(status, gin.H{
			"status":  "error",
			"message": message,
		})
	}
}

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
		handleError(c, status, "Employee not found", err)
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
		handleError(c, http.StatusInternalServerError, "Failed to retrieve employees", err)
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
		handleError(c, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	// Check if the associated user is active
	var user models.User
	if err := db.First(&user, "id = ?", employee.UserID).Error; err != nil {
		handleError(c, http.StatusBadRequest, "Associated user does not exist", err)
		return
	}

	if !user.IsActive {
		handleError(c, http.StatusBadRequest, "Cannot create employee for an inactive user", nil)
		return
	}

	// Validate the employee data
	if err := validateEmployee(db, &employee, true); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	employee.CreatedAt = time.Now()

	if err := db.Create(&employee).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create employee", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": employee})
}

// UpdateEmployee updates an existing employee
// UpdateEmployee updates an existing employee
func UpdateEmployee(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var employee models.Employee

	// Retrieve the employee by ID
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		status := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
			err = nil
		}
		handleError(c, status, "Employee not found", err)
		return
	}

	var input struct {
		Fullname   string `json:"fullname"`
		Phone      string `json:"phone"`
		PositionID int    `json:"position_id"`
		DivisionID int    `json:"division_id"`
		Avatar     string `json:"avatar"`
	}

	// Bind input JSON to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid input data", err)
		return
	}

	// Only update the fields that were provided and are different
	if input.Fullname != "" && input.Fullname != employee.Fullname {
		employee.Fullname = input.Fullname
	}
	if input.Phone != "" && input.Phone != employee.Phone {
		employee.Phone = input.Phone
	}
	if input.PositionID != 0 && input.PositionID != employee.PositionID {
		employee.PositionID = input.PositionID
	}
	if input.DivisionID != 0 && input.DivisionID != employee.DivisionID {
		employee.DivisionID = input.DivisionID
	}
	if input.Avatar != "" && input.Avatar != employee.Avatar {
		employee.Avatar = input.Avatar
	}

	// Validate only the fields that were changed
	if input.Fullname != "" || input.Phone != "" || input.PositionID != 0 || input.DivisionID != 0 {
		if err := validateEmployee(db, &employee, true); err != nil {
			handleError(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}

	// Set updated timestamp
	employee.UpdatedAt = time.Now()

	// Save the updated employee record
	if err := db.Save(&employee).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update employee", err)
		return
	}

	// Return the updated employee data
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": employee})
}

// DeleteEmployee deletes an employee by ID
func DeleteEmployee(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var employee models.Employee

	// Retrieve the employee by ID
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		handleError(c, http.StatusNotFound, "Employee not found", err)
		return
	}

	// Check if employee is active
	if employee.IsActive {
		handleError(c, http.StatusBadRequest, "Active employees cannot be deleted", nil)
		return
	}

	// Check if the employee has attendance data
	var attendanceCount int64
	if err := db.Model(&models.Attendance{}).Where("user_id = ?", employee.UserID).Count(&attendanceCount).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to check attendance", err)
		return
	}

	if attendanceCount > 0 {
		handleError(c, http.StatusBadRequest, "Employee has attendance data, cannot be deleted", nil)
		return
	}

	// Check if the associated user is active
	var user models.User
	if err := db.First(&user, "id = ?", employee.UserID).Error; err != nil {
		handleError(c, http.StatusBadRequest, "Associated user does not exist, cannot delete employee", err)
		return
	}

	if user.IsActive {
		handleError(c, http.StatusBadRequest, "Cannot delete employee with an active associated user", nil)
		return
	}

	// Proceed to delete the employee
	if err := db.Delete(&employee).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete employee", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Employee deleted successfully",
	})
}

// validateEmployee performs validation checks for employee data
func validateEmployee(db *gorm.DB, employee *models.Employee, isUpdate bool) error {
	// Validate Full Name
	if strings.TrimSpace(employee.Fullname) == "" {
		return errors.New("full name is required and cannot be empty")
	}

	// Validate User ID Uniqueness (excluding the current employee if it's an update)
	if err := db.Where("user_id = ? AND id != ?", employee.UserID, employee.ID).First(&models.Employee{}).Error; err == nil {
		return errors.New("user ID already associated with another employee")
	}

	// Validate Full Name Uniqueness (excluding the current employee if it's an update)
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

	// Skip Avatar validation if it's an update and the Avatar is not being updated
	if isUpdate && employee.Avatar == "" {
		return nil
	}

	// Further validation related to Avatar can be added here if needed

	return nil
}
