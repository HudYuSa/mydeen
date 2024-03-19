package controllers

import (
	"encoding/json"
	"log"

	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type WebSocketController interface {
	UpgradeConnection(ctx *gin.Context)
	HandleConnect(s *melody.Session)
	HandleDisconnect(s *melody.Session)
	HandleMessage(s *melody.Session, b []byte)
}

type webSocketController struct {
	QuestionController QuestionController
	LikeController     LikeController
	Melody             *melody.Melody
}

func NewWebSocketController(questionController QuestionController, likeController LikeController, m *melody.Melody) WebSocketController {
	return &webSocketController{
		QuestionController: questionController,
		LikeController:     likeController,
		Melody:             m,
	}
}

// UpgradeCConnection upgrades an HTTP connection to a WebSocket
func (wsc *webSocketController) UpgradeConnection(ctx *gin.Context) {
	// wsc.Melody.
	wsc.Melody.HandleRequest(ctx.Writer, ctx.Request.WithContext(ctx))
}

// HandleConnect handles new WebSocket Connections
func (wsc *webSocketController) HandleConnect(s *melody.Session) {
	log.Println("new connection")
	log.Println("connections: ", wsc.Melody.Len()+1)
}

// HandleDisconnect handles WebSocket Disconnections
func (wsc *webSocketController) HandleDisconnect(s *melody.Session) {
	log.Println("removed connection")
	log.Println("connections: ", wsc.Melody.Len()-1)
}

// HandleMessage handles incoming Websocket messages.
func (wsc *webSocketController) HandleMessage(s *melody.Session, b []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(b, &msg); err != nil {
		// handle decoding error
		// write back to the websocket connection session
		s.Write(dtos.EncodeJson(dtos.WebResponse{
			Message: err.Error(),
			Error:   true,
		}))
	}

	log.Println("new message")
	log.Println(msg)

	// Handle different message types
	switch msg["type"] {
	// questions message
	case string(dtos.CreateQuestionType):
		log.Println("entering create question type")
		wsc.QuestionController.CreateQuestion(s, b)

	case string(dtos.DeleteQuestionType):
		log.Println("entering delete question type")
		// example of using middleware
		// if middlewares.WSAuthenticateAdmin(s, dtos.Question) {
		// 	wsc.QuestionController.DeleteQuestion(s, b)
		// }
		wsc.QuestionController.DeleteQuestion(s, b)

	case string(dtos.EditQuestionType):
		log.Println("entering edit question type")
		wsc.QuestionController.EditQuestion(s, b)

	case string(dtos.AdminDeleteQuestionType):
		log.Println("entering admin delete question type")
		// not implemented

	case string(dtos.AdminEditQuestionType):
		log.Println("entering admin edit question type")
		// not implemented

		// likes message
	case string(dtos.ToggleLikeType):
		log.Println("entering toggle like type")
		wsc.LikeController.ToggleLike(s, b)
	}

}
