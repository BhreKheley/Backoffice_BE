package handlers

import (
	"absensi-app/models"
	"net/http"
	"strconv"
	"time"
	

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

// CheckIn handler to record clock-in
func CheckIn(c *gin.Context, db *sqlx.DB) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid", "error": err.Error()})
		return
	}

	attendance.ClockIn = time.Now()

	_, err := db.NamedExec(`INSERT INTO attendance (user_id, status_id, clock_in, clock_in_photo, latitude_clock_in, longitude_clock_in, description) 
        VALUES (:user_id, :status_id, :clock_in, :clock_in_photo, :latitude_clock_in, :longitude_clock_in, :description)`, &attendance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat check-in", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Check-in berhasil"})
}

// CheckOut handler to record clock-out
func CheckOut(c *gin.Context, db *sqlx.DB) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Input tidak valid", "error": err.Error()})
		return
	}

	attendance.ClockOut = time.Now()

	// Query untuk memastikan bahwa attendance record ada dan pengguna telah melakukan check-in
	_, err := db.NamedExec(`UPDATE attendance SET clock_out = :clock_out, clock_out_photo = :clock_out_photo, 
        latitude_clock_out = :latitude_clock_out, longitude_clock_out = :longitude_clock_out 
        WHERE id = :id AND user_id = :user_id AND clock_out IS NULL`, &attendance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mencatat check-out", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Check-out berhasil"})
}

// GetAttendanceByID retrieves attendance record by ID
func GetAttendanceByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")

	// Validasi ID
	attendanceID, err := strconv.Atoi(id)
	if err != nil || attendanceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID absensi tidak valid"})
		return
	}

	var attendance models.Attendance
	err = db.Get(&attendance, "SELECT * FROM attendance WHERE id = $1", attendanceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Gagal mengambil catatan absensi", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": attendance})
}

// GetAllAttendance retrieves all attendance records with pagination and optional date filtering
// GetAllAttendance retrieves all attendance records with pagination and optional date filtering
func GetAllAttendance(c *gin.Context, db *sqlx.DB) {
	page := c.Query("page")
	limit := c.Query("limit")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	var attendances []models.Attendance
	var query string
	var args []interface{}

	query = "SELECT * FROM attendance"
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

		query += " WHERE clock_in BETWEEN $1 AND $2"
		args = append(args, start, end)
	}
	query += " ORDER BY clock_in DESC LIMIT $3 OFFSET $4"
	args = append(args, limitNum, offset)

	err = db.Select(&attendances, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil catatan absensi"})
		return
	}

	var totalRecords int
	err = db.Get(&totalRecords, "SELECT COUNT(*) FROM attendance")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghitung total catatan"})
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
			TotalRecords: totalRecords,
			Page:         pageNum,
			Limit:        limitNum,
		},
		Attendances: formattedAttendances,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data retrieved successfully", "data": response})
}
