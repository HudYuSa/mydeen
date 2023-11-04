package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.RouterGroup) {
	// Routes
	master := NewMasterRoutes(controllers.Master)
	common := NewCommonRoutes(controllers.Common)
	admin := NewAdminRoutes(controllers.Admin)
	event := NewEventRoutes(controllers.Event)

	// setup routes
	master.SetupRoutes(router)
	common.SetupRoutes(router)
	admin.SetupRoutes(router)
	event.SetupRoutes(router)
}
