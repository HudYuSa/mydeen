package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/gin-gonic/gin"
)

// routes is only created for an http endpoint
// for websocket endpoint it only needs the controller

func InitializeRoutes(router *gin.RouterGroup) {
	// Routes
	master := NewMasterRoutes(controllers.Master)
	common := NewCommonRoutes(controllers.Common)
	admin := NewAdminRoutes(controllers.Admin)
	event := NewEventRoutes(controllers.Event)
	question := NewQuestionRoutes(controllers.Question)
	webSocket := NewWebSocketController(controllers.WebSocket)

	// setup routes
	master.SetupRoutes(router)
	common.SetupRoutes(router)
	admin.SetupRoutes(router)
	event.SetupRoutes(router)
	question.SetupRoutes(router)
	webSocket.SetupRoutes(router)
}
