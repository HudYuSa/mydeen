package models

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	QuestionID uuid.UUID
	EventID    uuid.UUID
	UserID     uuid.UUID
	Content    string
	Starred    bool
	Approved   bool
	Answered   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Event      Event `gorm:"foreignKey:EventID"`
	User       User  `gorm:"foreignKey:UserID"`
}
