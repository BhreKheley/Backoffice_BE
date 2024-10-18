package models

import (
	"time"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Fullname  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
