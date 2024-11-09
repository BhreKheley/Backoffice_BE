package handlers

import (
	"absensi-app/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetEmployee retrieves an employee by ID
func GetEmployee(c *gin.Context, db *gorm.DB) {
	var employee models.Employee
	id := c.Param("id")

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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   employee,
	})
}

// GetAllEmployees retrieves all employees
func GetAllEmployees(c *gin.Context, db *gorm.DB) {
	var employees []struct {
		models.Employee
		PositionName string `json:"position_name" gorm:"column:position_name"`
		DivisionName string `json:"division_name" gorm:"column:division_name"`
	}

	err := db.Table(`"employee" AS e`).
		Select("e.id, e.user_id, e.full_name, e.phone, e.position_id, e.division_id, e.is_active, e.created_at, e.updated_at, p.position_name, d.division_name").
		Joins("LEFT JOIN position AS p ON e.position_id = p.id").
		Joins("LEFT JOIN division AS d ON e.division_id = d.id").
		Scan(&employees).Error

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve employees"})
		return
	} else {
		fmt.Printf("Employees: %+v\n", employees)
	}

	// if err := db.Find(&employees).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  "error",
	// 		"message": "Failed to retrieve employees",
	// 		"error":   err.Error(),
	// 	})
	// 	return
	// }

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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var count int64
	db.Model(&models.Employee{}).Where("user_id = ?", employee.UserID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User ini sudah terkait dengan employee lain",
		})
		return
	}

	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	if err := db.Create(&employee).Error; err != nil {
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
func UpdateEmployee(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var employee models.Employee

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

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	employee.UpdatedAt = time.Now()

	if err := db.Save(&employee).Error; err != nil {
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
func DeleteEmployee(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var employee models.Employee

	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
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

	var attendanceCount int64
	db.Model(&models.Attendance{}).Where("user_id = ?", employee.UserID).Count(&attendanceCount)
	if attendanceCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Employee memiliki data absensi, tidak dapat dihapus",
		})
		return
	}

	var userIsActive bool
	if err := db.Model(&models.User{}).Select("is_active").Where("id = ?", employee.UserID).Scan(&userIsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memeriksa status user",
			"error":   err.Error(),
		})
		return
	}

	if userIsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User terkait masih aktif, employee tidak dapat dihapus",
		})
		return
	}

	if err := db.Delete(&employee).Error; err != nil {
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
