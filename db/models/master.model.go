package models

import (
	"time"

	"github.com/google/uuid"
)

type Master struct {
	MasterID  uuid.UUID
	Email     string
	Password  string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
