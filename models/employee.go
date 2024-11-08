package models

import "time"

// Employee represents the employee model
type Employee struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id"`
	Fullname   string    `db:"full_name" json:"full_name"`
	Phone      string    `db:"phone" json:"phone"`
	PositionID int       `db:"position_id" json:"position_id"`
	DivisionID int       `db:"division_id" json:"division_id"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
