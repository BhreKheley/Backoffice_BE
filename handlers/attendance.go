package handlers

import (
	"absensi-app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AttendanceResponse represents the response format for the attendance API
type AttendanceResponse struct {
	Metadata    Metadata              `json:"metadata"`
	Attendances []AttendanceFormatted `json:"attendances"`
}

// Metadata contains pagination information
type Metadata struct {
	TotalRecords int `json:"total_records"`
	Page         int `json:"page"`
	Limit        int `json:"limit"`
}

// AttendanceFormatted represents a formatted attendance record without redundant fields
type AttendanceFormatted struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	StatusID          int       `json:"status_id"`
	ClockIn           time.Time `json:"clock_in"`
	ClockInPhoto      string    `json:"clock_in_photo"`
	LatitudeClockIn   float64   `json:"latitude_clock_in"`
	LongitudeClockIn  float64   `json:"longitude_clock_in"`
	ClockOut          time.Time `json:"clock_out"`
	ClockOutPhoto     string    `json:"clock_out_photo"`
	LatitudeClockOut  float64   `json:"latitude_clock_out"`
	LongitudeClockOut float64   `json:"longitude_clock_out"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CheckIn handler untuk mencatat clock-in
func CheckIn(c *gin.Context, db *gorm.DB) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid", "error": err.Error()})
		return
	}

	attendance.ClockIn = time.Now()

	if err := db.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat check-in", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Check-in berhasil"})
}

// CheckOut handler untuk mencatat clock-out
func CheckOut(c *gin.Context, db *gorm.DB) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid", "error": err.Error()})
		return
	}

	attendance.ClockOut = time.Now()

	// Update attendance untuk user yang belum melakukan check-out
	result := db.Model(&models.Attendance{}).
		Where("id = ? AND user_id = ? AND clock_out IS NULL", attendance.ID, attendance.UserID).
		Updates(map[string]interface{}{
			"clock_out":          attendance.ClockOut,
			"clock_out_photo":    attendance.ClockOutPhoto,
			"latitude_clock_out": attendance.LatitudeClockOut,
			"longitude_clock_out": attendance.LongitudeClockOut,
		})

	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat check-out", "error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Check-out berhasil"})
}

// GetAttendanceByID handler untuk mengambil catatan absensi berdasarkan ID
func GetAttendanceByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	attendanceID, err := strconv.Atoi(id)
	if err != nil || attendanceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID absensi tidak valid"})
		return
	}

	var attendance models.Attendance
	if err := db.First(&attendance, attendanceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Catatan absensi tidak ditemukan", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": attendance})
}

// GetAllAttendance handler untuk mengambil semua catatan absensi dengan pagination dan filter tanggal
func GetAllAttendance(c *gin.Context, db *gorm.DB) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	var attendances []models.Attendance
	query := db.Order("clock_in DESC")

	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format tanggal awal tidak valid"})
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format tanggal akhir tidak valid"})
			return
		}
		query = query.Where("clock_in BETWEEN ? AND ?", start, end)
	}

	totalRecords := int64(0)
	query.Model(&models.Attendance{}).Count(&totalRecords)

	query = query.Limit(limit).Offset((page - 1) * limit)
	if err := query.Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil catatan absensi"})
		return
	}

	var formattedAttendances []AttendanceFormatted
	for _, a := range attendances {
		formattedAttendances = append(formattedAttendances, AttendanceFormatted{
			ID:                a.ID,
			UserID:            a.UserID,
			StatusID:          a.StatusID,
			ClockIn:           a.ClockIn,
			ClockInPhoto:      a.ClockInPhoto,
			LatitudeClockIn:   a.LatitudeClockIn,
			LongitudeClockIn:  a.LongitudeClockIn,
			ClockOut:          a.ClockOut,
			ClockOutPhoto:     a.ClockOutPhoto,
			LatitudeClockOut:  a.LatitudeClockOut,
			LongitudeClockOut: a.LongitudeClockOut,
			Description:       a.Description,
			CreatedAt:         a.CreatedAt,
			UpdatedAt:         a.UpdatedAt,
		})
	}

	response := AttendanceResponse{
		Metadata: Metadata{
			TotalRecords: int(totalRecords),
			Page:         page,
			Limit:        limit,
		},
		Attendances: formattedAttendances,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data retrieved successfully", "data": response})
}
