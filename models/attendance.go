package models

import (
	"time"
	"gorm.io/gorm"
)


// Attendance represents the attendance record of an employee
type Attendance struct {
	gorm.Model
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	StatusID   int       `json:"status_id"`
	ClockIn    time.Time `json:"clock_in"`
	ClockInPhoto string    `json:"clock_in_photo"`
	LatitudeClockIn   float64   `json:"latitude_clock_in"`
	LongitudeClockIn  float64   `json:"longitude_clock_in"`
	ClockOut   time.Time `json:"clock_out"`
	ClockOutPhoto string    `json:"clock_out_photo"`
	LatitudeClockOut   float64   `json:"latitude_clock_out"`
	LongitudeClockOut  float64   `json:"longitude_clock_out"`
	Description string    `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
