package models

import (
	"gorm.io/gorm"
)

// Role represents the role of a user in the system
type Role struct {
	gorm.Model
	ID   int    `json:"id"`
	Rolename string `json:"role_name"`
	Code string `json:"code"`
}
