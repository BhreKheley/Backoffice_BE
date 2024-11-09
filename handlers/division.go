package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"absensi-app/models"
)

// GetAllDivisions retrieves all divisions
func GetAllDivisions(c *gin.Context, db *gorm.DB) {
	var divisions []models.Division
	if err := db.Find(&divisions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, divisions)
}

// GetDivisionByID retrieves a division by ID
func GetDivisionByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var division models.Division
	if err := db.First(&division, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}
	c.JSON(http.StatusOK, division)
}

// CreateDivision creates a new division
func CreateDivision(c *gin.Context, db *gorm.DB) {
	var division models.Division
	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := db.Create(&division).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, division)
}

// UpdateDivision updates an existing division
func UpdateDivision(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var division models.Division
	if err := db.First(&division, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}

	if err := c.ShouldBindJSON(&division); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	division.ID = int(id) // Ensure the ID stays consistent
	if err := db.Save(&division).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, division)
}

// DeleteDivision deletes a division by ID
func DeleteDivision(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.Delete(&models.Division{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
