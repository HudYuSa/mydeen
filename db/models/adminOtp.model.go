package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminOtp struct {
	AdminOtpID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	AdminID    uuid.UUID `gorm:"not null"`
	Code       string    `gorm:"not null"`
	ExpireDate time.Time `gorm:"not null"`
	Used       bool      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null"`
	Admin      Admin     `gorm:"foreignKey:AdminID;references:AdminID"`
}
