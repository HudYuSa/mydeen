package controllers

import (
	"github.com/HudYuSa/mydeen/internal/connection"
)

var (
	Master MasterController
	Common CommonController
	Admin  AdminController
	Event  EventController
)

func InitializeControllers() {
	Common = NewCommonController(connection.DB)
	Master = NewMasterController(connection.DB)
	Admin = NewAdminController(connection.DB)
	Event = NewEventController(connection.DB)
}
