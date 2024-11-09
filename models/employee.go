package models

import (
	"time"
)

// Employee represents the employee model
type Employee struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"not null" json:"user_id"`
	Fullname   string    `gorm:"column:full_name" json:"full_name"`
	Phone      string    `gorm:"type:varchar(20)" json:"phone"`
	PositionID int       `gorm:"not null" json:"position_id"`
	DivisionID int       `gorm:"not null" json:"division_id"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Custom Table Name for RolePermission
func (Employee) TableName() string {
	return "employee"
}
