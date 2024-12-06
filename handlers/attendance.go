package handlers

import (
	"absensi-app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckIn handles the check-in functionality
func CheckIn(c *gin.Context, db *gorm.DB) {
	var attendance models.Attendance

	// Bind JSON input
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Validate User ID
	var userExists bool
	if err := db.Model(&models.User{}).
		Select("count(*) > 0").
		Where("id = ?", attendance.UserID).
		Find(&userExists).Error; err != nil || !userExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid User ID",
		})
		return
	}

	// Check if the user is already clocked in today
	today := time.Now().Format("2006-01-02")
	var existingAttendance models.Attendance
	if err := db.Where("user_id = ? AND DATE(clock_in) = ?", attendance.UserID, today).
		First(&existingAttendance).Error; err == nil && existingAttendance.IsClockedIn {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User already checked in today",
		})
		return
	}

	// Record Clock-In
	attendance.ClockIn = time.Now()
	attendance.IsClockedIn = true

	if err := db.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to record check-in",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Check-in successful",
		"data":    attendance,
	})
}

// CheckOut handles the check-out functionality
func CheckOut(c *gin.Context, db *gorm.DB) {
	var input struct {
		UserID            int     `json:"user_id" binding:"required"`
		StatusID          int     `json:"status_id"`
		ClockOutPhoto     string  `json:"clock_out_photo"`
		LatitudeClockOut  float64 `json:"latitude_clock_out"`
		LongitudeClockOut float64 `json:"longitude_clock_out"`
		Description       string  `json:"description"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Validasi keberadaan status ID (jika diberikan)
	if input.StatusID != 0 {
		var statusExists bool
		if err := db.Model(&models.Status{}).
			Select("count(*) > 0").
			Where("id = ?", input.StatusID).
			Find(&statusExists).Error; err != nil || !statusExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid Status ID",
			})
			return
		}
	}

	// Temukan data kehadiran yang sedang aktif (Clock-In tanpa Clock-Out)
	var attendance models.Attendance

	// Ambil tanggal hari ini
	today := time.Now().Format("2006-01-02")

	// Periksa catatan yang sesuai dengan user_id dan tanggal hari ini
	if err := db.Where("user_id = ? AND DATE(clock_in) = ? AND is_clocked_in = true AND is_clocked_out = false", input.UserID, today).
		First(&attendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No attendance record for today found to check-out",
			"error":   err.Error(),
		})
		return
	}

	// Update Clock-Out details
	now := time.Now()
	attendance.ClockOut = &now
	attendance.IsClockedOut = true
	attendance.ClockOutPhoto = input.ClockOutPhoto
	attendance.LatitudeClockOut = input.LatitudeClockOut
	attendance.LongitudeClockOut = input.LongitudeClockOut
	attendance.Description = input.Description
	if input.StatusID != 0 {
		attendance.StatusID = input.StatusID
	}

	if err := db.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to record check-out",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Check-out successful",
		"data":    attendance,
	})
}

// GetAttendanceByUserID retrieves all attendance records of a specific user with pagination
func GetAttendanceByUserID(c *gin.Context, db *gorm.DB) {
	// Ambil user_id dari parameter URL
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid User ID",
		})
		return
	}

	// Ambil query untuk pagination (page dan limit)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1 // Set default jika query tidak valid
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10 // Set default jika query tidak valid
	}

	offset := (page - 1) * limit

	// Validasi apakah user dengan ID tertentu ada
	var userExists bool
	if err := db.Model(&models.User{}).
		Select("count(*) > 0").
		Where("id = ?", userID).
		Find(&userExists).Error; err != nil || !userExists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	// Ambil semua data kehadiran user dengan pagination
	var attendances []models.Attendance
	var totalRecords int64

	// Hitung total data kehadiran
	if err := db.Model(&models.Attendance{}).
		Where("user_id = ?", userID).
		Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to count attendance records",
			"error":   err.Error(),
		})
		return
	}

	// Ambil data dengan limit dan offset
	if err := db.Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("clock_in DESC").
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve attendance records",
			"error":   err.Error(),
		})
		return
	}

	// Periksa apakah data kosong
	if len(attendances) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No attendance records found for this user",
			"data":    []models.Attendance{},
		})
		return
	}

	// Kirim data response dalam bentuk JSON
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"message":       "Attendance records retrieved successfully",
		"total_records": totalRecords,
		"current_page":  page,
		"total_pages":   (totalRecords + int64(limit) - 1) / int64(limit), // Hitung total halaman
		"data":          attendances,
	})
}

// GetUserAttendanceStatus retrieves the attendance status of a specific user
func GetUserAttendanceStatus(c *gin.Context, db *gorm.DB) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid User ID",
		})
		return
	}

	// Get today's date
	today := time.Now().Format("2006-01-02")

	// Find the attendance record for today
	var attendance models.Attendance
	err = db.Where("user_id = ? AND DATE(clock_in) = ?", userID, today).
		First(&attendance).Error

	message := "Kamu belum Clock In hari ini"
	if err == nil {
		if attendance.IsClockedOut {
			message = "Tugas Kamu hari ini sudah selesai"
		} else if attendance.IsClockedIn {
			message = "Kamu sudah Clock In hari ini dan belum Clock Out"
		}
	} else if err != gorm.ErrRecordNotFound {
		// Handle database error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve attendance record",
			"error":   err.Error(),
		})
		return
	}

	// Return the user's attendance status
	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"message":         message,
		"user_id":         userID,
		"attendance_code": attendance.StatusID,
		"attendance_status": gin.H{
			"is_clocked_in":  attendance.IsClockedIn,
			"is_clocked_out": attendance.IsClockedOut,
		},
		"attendance_record": gin.H{
			"clock_in":            attendance.ClockIn,
			"clock_in_photo":      attendance.ClockInPhoto,
			"latitude_clock_in":   attendance.LatitudeClockIn,
			"longitude_clock_in":  attendance.LongitudeClockIn,
			"clock_out":           attendance.ClockOut,
			"clock_out_photo":     attendance.ClockOutPhoto,
			"latitude_clock_out":  attendance.LatitudeClockOut,
			"longitude_clock_out": attendance.LongitudeClockOut,
			"description":         attendance.Description,
		},
	})
}

// GetAttendanceByID retrieves an attendance record by its ID
func GetAttendanceByID(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid ID"})
		return
	}

	var attendance models.Attendance
	if err := db.First(&attendance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Attendance record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": attendance})
}

// GetAllAttendance retrieves all attendance records with pagination
func GetAllAttendance(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var attendances []models.Attendance
	var totalRecords int64

	if err := db.Model(&models.Attendance{}).Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to count attendance records"})
		return
	}

	if err := db.Limit(limit).Offset(offset).Order("clock_in DESC").Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve attendance records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   attendances,
	})
}
