package models

// Role represents the role of a user in the system
type Role struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName string `gorm:"type:varchar(255);unique;not null" json:"role_name"`
	Code     string `gorm:"type:varchar(100);unique;not null" json:"code"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}

// Custom Table Name for Role
func (Role) TableName() string {
	return "role"
}
