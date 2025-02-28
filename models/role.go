package models

type Role struct {
	ID         int          `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName   string       `gorm:"type:varchar(255);unique;not null" json:"role_name"`
	Code       string       `gorm:"type:varchar(100);unique;not null" json:"code"`
	IsActive   bool         `gorm:"not null" json:"is_active"`
	Permissions []Permission `gorm:"many2many:role_permission" json:"permissions"` // Tambahkan ini
}

// Custom Table Name for Role
func (Role) TableName() string {
	return "role"
}
