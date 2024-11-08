package models

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           int `db:"id" json:"id"` // Ganti gorm menjadi db
	RoleID       int `db:"role_id" json:"role_id"` // Ganti gorm menjadi db
	PermissionID int `db:"permission_id" json:"permission_id"` // Ganti gorm menjadi db
}

// Custom Table Name for RolePermission
func (RolePermission) TableName() string {
	return "role_permission"
}
