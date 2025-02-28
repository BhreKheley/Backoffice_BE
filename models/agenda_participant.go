package models

import (
	"time"
)

type AgendaParticipant struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AgendaID  uint      `gorm:"not null" json:"agenda_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Role      string    `gorm:"type:varchar(255);default:'Peserta'" json:"role"`
	AddedBy   uint      `gorm:"not null" json:"added_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi
	Agenda Agenda `gorm:"foreignKey:AgendaID" json:"agenda"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
	Adder  User   `gorm:"foreignKey:AddedBy" json:"adder"`
}

func (AgendaParticipant) TableName() string {
	return "agenda_participants"
}
