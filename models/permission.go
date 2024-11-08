package models

// Permission represents the permission in the system
type Permission struct {
	ID             int    `db:"id" json:"id"` // Ganti gorm menjadi db
	PermissionName string `db:"permission_name" json:"permission_name"` // Ganti gorm menjadi db
	Code           string `db:"code" json:"code"` // Ganti gorm menjadi db
}

// Custom Table Name for Permission
func (Permission) TableName() string {
	return "permission"
}
