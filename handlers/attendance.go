package handlers

import (
	"absensi-app/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

func CheckIn(c *gin.Context, db *sqlx.DB) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance.ClockIn = time.Now()

	_, err := db.NamedExec(`INSERT INTO attendance (user_id, status_id, clock_in, clock_in_photo, latitude_clock_in, longitude_clock_in, description) 
        VALUES (:user_id, :status_id, :clock_in, :clock_in_photo, :latitude_clock_in, :longitude_clock_in, :description)`, &attendance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check-in successful!"})
}

func CheckOut(c *gin.Context, db *sqlx.DB) {
	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	attendance.ClockOut = time.Now()

	_, err := db.NamedExec(`UPDATE attendance SET clock_out = :clock_out, clock_out_photo = :clock_out_photo, 
        latitude_clock_out = :latitude_clock_out, longitude_clock_out = :longitude_clock_out 
        WHERE id = :id AND user_id = :user_id`, &attendance)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check-out successful!"})
}

func GetAttendanceByUser(c *gin.Context, db *sqlx.DB) {
	userID := c.Param("user_id")

	var attendances []models.Attendance
	err := db.Select(&attendances, "SELECT * FROM attendance WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendances)
}

