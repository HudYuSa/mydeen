package models

import (
	"net"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID
	Username  *string
	Email     *string
	IpAddress *net.IP
}
