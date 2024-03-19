package models

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	QuestionID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	EventID    uuid.UUID `gorm:"not null"`
	UserID     uuid.UUID `gorm:"not null"`
	Username   string    `gorm:"not null"`
	Content    string    `gorm:"not null"`
	Starred    bool      `gorm:"not null"`
	Approved   bool      `gorm:"not null"`
	Answered   bool      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	Event      Event     `gorm:"foreignKey:EventID;references:EventID"`
	Likes      []Like    `gorm:"references:QuestionID"`
}
