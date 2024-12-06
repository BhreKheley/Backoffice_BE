package models

import (
	"time"
)

type Position struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	PositionName string    `gorm:"type:varchar(255);not null" json:"position_name"`
	DivisionID   int       `gorm:"not null" json:"division_id"`
	IsActive     bool      `gorm:"not null" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Custom Table Name for RolePermission
func (Position) TableName() string {
	return "position"
}
