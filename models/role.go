package models

// Role represents the role of a user in the system
type Role struct {
	ID       int    `db:"id" json:"id"`
	RoleName string `db:"role_name" json:"role_name"` // Pastikan tag db di sini
	Code     string `db:"code" json:"code"`
}

// Custom Table Name for Role
func (Role) TableName() string {
	return "role"
}
