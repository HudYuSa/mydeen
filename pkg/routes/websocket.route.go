package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type WebSocketRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type webSocketRoutes struct {
	WebSocketController controllers.WebSocketController
}

func NewWebSocketController(webSocketController controllers.WebSocketController) WebSocketRoutes {
	return &webSocketRoutes{
		WebSocketController: webSocketController,
	}
}

func (wsr *webSocketRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/ws")

	router.GET("", controllers.WebSocket.UpgradeConnection)
}
