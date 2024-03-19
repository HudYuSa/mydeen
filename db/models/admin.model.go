package models

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	AdminID      uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	InvitationID uuid.UUID  `gorm:"not null"`
	Username     string     `gorm:"not null"`
	Email        string     `gorm:"uniqueIndex;not null"`
	Password     string     `gorm:"not null"`
	AdminCode    string     `gorm:"not null"`
	Enable2fa    bool       `gorm:"not null;column:enable_2fa"`
	CreatedAt    time.Time  `gorm:"not null"`
	UpdatedAt    time.Time  `gorm:"not null"`
	Invitation   Invitation `gorm:"foreignKey:InvitationID;references:InvitationID"`
}
