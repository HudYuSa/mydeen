package models

import (
	"time"

	"github.com/google/uuid"
)

type MasterOtp struct {
	MasterOtpID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	MasterID    uuid.UUID `gorm:"not null"`
	Code        string    `gorm:"not null"`
	ExpireDate  time.Time `gorm:"not null"`
	Used        bool      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	Master      Master    `gorm:"foreignKey:MasterID;references:MasterID"`
}
