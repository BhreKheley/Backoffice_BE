package models

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           int `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       int `gorm:"not null" json:"role_id"`
	PermissionID int `gorm:"not null" json:"permission_id"`
}

// Custom Table Name for RolePermission
func (RolePermission) TableName() string {
	return "role_permission"
}
