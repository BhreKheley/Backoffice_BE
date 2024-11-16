package models

// Role represents the role of a user in the system
type Role struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName string `gorm:"type:varchar(255);not null" json:"role_name"`
	Code     string `gorm:"type:varchar(100);not null" json:"code"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// Custom Table Name for Role
func (Role) TableName() string {
	return "role"
}
