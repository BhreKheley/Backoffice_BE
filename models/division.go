package models

import (
	"time"
)

type Division struct {
	ID             int       `db:"id" json:"id"`
	DivisionName   string    `db:"division_name" json:"division_name"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
