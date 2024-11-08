package models

import (
	"time"
)

type Position struct {
	ID           int       `db:"id" json:"id"`
	PositionName string    `db:"position_name" json:"position_name"`
	DivisionID   int       `db:"division_id" json:"division_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
