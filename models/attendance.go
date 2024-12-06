package models

import (
	"time"
)

// Attendance represents the attendance record of an employee
type Attendance struct {
	ID                int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            int        `gorm:"not null" json:"user_id"`
	StatusID          int        `gorm:"not null" json:"status_id"`
	ClockIn           time.Time  `gorm:"type:timestamp" json:"clock_in"`
	ClockInPhoto      string     `gorm:"type:text" json:"clock_in_photo"`
	LatitudeClockIn   float64    `gorm:"type:float" json:"latitude_clock_in"`
	LongitudeClockIn  float64    `gorm:"type:float" json:"longitude_clock_in"`
	IsClockedIn       bool       `gorm:"not null;default:false" json:"is_clocked_in"`
	ClockOut          *time.Time `gorm:"type:timestamp" json:"clock_out"` // Pointer for nullable field
	ClockOutPhoto     string     `gorm:"type:text" json:"clock_out_photo"`
	LatitudeClockOut  float64    `gorm:"type:float" json:"latitude_clock_out"`
	LongitudeClockOut float64    `gorm:"type:float" json:"longitude_clock_out"`
	IsClockedOut      bool       `gorm:"not null;default:false" json:"is_clocked_out"`
	Description       string     `gorm:"type:varchar(255)" json:"description"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// Custom Table Name for Attendance
func (Attendance) TableName() string {
	return "attendance"
}
