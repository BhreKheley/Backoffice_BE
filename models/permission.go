package models

type Permission struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"id"`
	PermissionName string `gorm:"type:varchar(255);not null" json:"permission_name"`
	Code           string `gorm:"type:varchar(100);not null" json:"code"`
	IsActive       bool   `gorm:"not null" json:"is_active"`
	Roles          []Role `gorm:"many2many:role_permission" json:"roles"` // Tambahkan ini
}

// Custom Table Name for Permission
func (Permission) TableName() string {
	return "permission"
}
