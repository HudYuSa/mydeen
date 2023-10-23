package models

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	AdminID      uuid.UUID
	InvitationID uuid.UUID
	Username     string
	Email        string
	Password     string
	Enable2fa    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Invitation   Invitation `gorm:"foreignKey:InvitationID"`
}
