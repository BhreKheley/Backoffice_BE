package handlers

import (
	"net/http"
	"strconv"
	"time"

	"absensi-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllPositions retrieves all positions
func GetAllPositions(c *gin.Context, db *gorm.DB) {
	var positions []models.Position
	if err := db.Find(&positions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, positions)
}

// GetPositionByID retrieves a position by ID
func GetPositionByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var position models.Position
	if err := db.First(&position, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, position)
}

// CreatePosition creates a new position
func CreatePosition(c *gin.Context, db *gorm.DB) {
	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	position.CreatedAt = time.Now()
	position.UpdatedAt = time.Now()

	if err := db.Create(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, position)
}

// UpdatePosition updates an existing position
func UpdatePosition(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var position models.Position
	if err := db.First(&position, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	position.ID = id
	position.UpdatedAt = time.Now()

	if err := db.Save(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, position)
}

// DeletePosition deletes a position by ID
func DeletePosition(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.Delete(&models.Position{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
