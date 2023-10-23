package models

import (
	"time"

	"github.com/google/uuid"
)

type MasterOtp struct {
	MasterOtpID uuid.UUID
	MasterID    uuid.UUID
	Code        string
	ExpireDate  time.Time
	Used        bool
	CreatedAt   time.Time
	Master      Master `gorm:"foreignKey:MasterID"`
}
