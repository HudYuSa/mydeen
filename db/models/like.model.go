package models

import "github.com/google/uuid"

type Like struct {
	LikeID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	QuestionID uuid.UUID `gorm:"not null"`
	UserID     uuid.UUID `gorm:"not null"`
}
