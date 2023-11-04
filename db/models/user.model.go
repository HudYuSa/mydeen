package models

import (
	"net"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  *string
	Email     *string
	IpAddress *net.IP
}
