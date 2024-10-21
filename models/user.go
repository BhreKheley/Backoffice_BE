package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int    `json:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	RoleID   int    `json:"role_id" db:"role_id"`
	IsActive bool   `json:"is_active" db:"is_active"`
}
