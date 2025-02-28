package models

import (
	"time"
)

type Agenda struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
	StartTime   string    `gorm:"type:varchar(10);not null" json:"start_time"`
	EndTime     string    `gorm:"type:varchar(10);not null" json:"end_time"`
	Location    string    `gorm:"type:varchar(255)" json:"location"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Creator User `gorm:"foreignKey:CreatedBy;references:id" json:"creator"`
}

func (Agenda) TableName() string {
	return "agendas"
}
