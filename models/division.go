package models

import (
	"time"
)

// Division represents a division in the organization
type Division struct {
	ID           int      `gorm:"primaryKey;autoIncrement" json:"id"`
	DivisionName string    `gorm:"type:varchar(255);not null" json:"division_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Custom Table Name for RolePermission
func (Division) TableName() string {
	return "division"
}