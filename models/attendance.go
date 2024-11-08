package models

import (
	"time"
	"gorm.io/gorm"
)

// Attendance represents the attendance record of an employee
type Attendance struct {
	gorm.Model
	ID                  int       `json:"id" db:"id"`
	UserID              int       `json:"user_id" db:"user_id"`
	StatusID            int       `json:"status_id" db:"status_id"`
	ClockIn             time.Time `json:"clock_in" db:"clock_in"`
	ClockInPhoto        string    `json:"clock_in_photo" db:"clock_in_photo"`
	LatitudeClockIn     float64   `json:"latitude_clock_in" db:"latitude_clock_in"`
	LongitudeClockIn    float64   `json:"longitude_clock_in" db:"longitude_clock_in"`
	ClockOut            time.Time `json:"clock_out" db:"clock_out"`
	ClockOutPhoto       string    `json:"clock_out_photo" db:"clock_out_photo"`
	LatitudeClockOut    float64   `json:"latitude_clock_out" db:"latitude_clock_out"`
	LongitudeClockOut   float64   `json:"longitude_clock_out" db:"longitude_clock_out"`
	Description         string    `json:"description" db:"description"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
