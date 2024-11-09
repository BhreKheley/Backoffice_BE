package models

import (
)

// Status represents the status of attendance
type Status struct {
	ID        int   `gorm:"primaryKey;autoIncrement" json:"id"`
	Statusname string `gorm:"type:varchar(255);not null" json:"status_name"`
	Code      string `gorm:"type:varchar(100);not null" json:"code"`
}

// TableName method untuk men-override nama tabel
func (Status) TableName() string {
	return "status"
}
