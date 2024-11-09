package models

// Permission represents the permission in the system
type Permission struct {
	ID             int   `gorm:"primaryKey;autoIncrement" json:"id"`
	PermissionName string `gorm:"type:varchar(255);not null" json:"permission_name"`
	Code           string `gorm:"type:varchar(100);not null" json:"code"`
}

// Custom Table Name for Permission
func (Permission) TableName() string {
	return "permission"
}
