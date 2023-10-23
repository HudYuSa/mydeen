package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminOtp struct {
	AdminOtpID uuid.UUID
	AdminID    uuid.UUID
	Code       string
	ExpireDate time.Time
	Used       bool
	CreatedAt  time.Time
	Admin      Admin `gorm:"foreignKey:AdminID"`
}
