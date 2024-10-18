package models

import (
	"gorm.io/gorm"
)

// Status represents the status of attendance
type Status struct {
	gorm.Model
	ID   int    `json:"id"`
	Statusname string `json:"status_name"`
	Code string `json:"code"`
}
