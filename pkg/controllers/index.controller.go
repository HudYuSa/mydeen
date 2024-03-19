package controllers

import (
	"github.com/HudYuSa/mydeen/internal/connection"
	"github.com/olahol/melody"
)

var (
	Master    MasterController
	Common    CommonController
	Admin     AdminController
	Event     EventController
	Question  QuestionController
	Like      LikeController
	WebSocket WebSocketController
)

func InitializeControllers(melody *melody.Melody) {
	Common = NewCommonController(connection.DB)
	Master = NewMasterController(connection.DB)
	Admin = NewAdminController(connection.DB)
	Event = NewEventController(connection.DB)
	Question = NewQuestionController(connection.DB, melody)
	Like = NewLikeController(connection.DB, melody)
	WebSocket = NewWebSocketController(Question, Like, melody)
}
