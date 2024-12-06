package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"absensi-app/models"
)

// GetAllStatus mendapatkan semua status
func GetAllStatus(c *gin.Context, db *gorm.DB) {
	var statuses []models.Status
	if err := db.Find(&statuses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch statuses", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": statuses})
}

// GetStatusByID mendapatkan status berdasarkan ID
func GetStatusByID(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ID"})
		return
	}

	var status models.Status
	if err := db.First(&status, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Status not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": status})
}

// CreateStatus membuat status baru
func CreateStatus(c *gin.Context, db *gorm.DB) {
	var newStatus models.Status
	if err := c.ShouldBindJSON(&newStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	// Cek apakah status_name atau code sudah ada
	var existingStatus models.Status
	if err := db.Where("statusname = ?", newStatus.Statusname).Or("code = ?", newStatus.Code).First(&existingStatus).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Status name or code already exists"})
		return
	}

	if err := db.Create(&newStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create status", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newStatus})
}

// UpdateStatus memperbarui status berdasarkan ID
func UpdateStatus(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ID"})
		return
	}

	var status models.Status
	if err := db.First(&status, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Status not found", "error": err.Error()})
		return
	}

	var input models.Status
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input", "error": err.Error()})
		return
	}

	// Cek apakah status_name atau code sudah ada (kecuali milik record ini sendiri)
	var existingStatus models.Status
	if err := db.Where("(statusname = ? OR code = ?) AND id != ?", input.Statusname, input.Code, id).First(&existingStatus).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Status name or code already exists"})
		return
	}

	// Update status
	status.Statusname = input.Statusname
	status.Code = input.Code
	status.IsActive = input.IsActive

	if err := db.Save(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update status", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": status})
}

// DeleteStatus menghapus status berdasarkan ID
func DeleteStatus(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ID"})
		return
	}

	var status models.Status
	if err := db.First(&status, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Status not found", "error": err.Error()})
		return
	}

	if err := db.Delete(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete status", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Status deleted successfully"})
}
