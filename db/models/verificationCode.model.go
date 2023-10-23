package models

import (
	"time"

	"github.com/google/uuid"
)

type VerificationCode struct {
	VerificationCodeID uuid.UUID
	MasterID           uuid.UUID
	Code               string
	ExpireDate         time.Time
	Used               bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Master             Master `gorm:"foreignKey:MasterID"`
}
