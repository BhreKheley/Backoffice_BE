package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
	IsActive bool   `json:"is_active"`
}