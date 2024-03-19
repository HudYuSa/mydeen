package dtos

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
)

// this package is where u put all data transfer object representation
// and all its function

// this is for the message group of which the message are send to
type WebSocketGroup string

const (
	Question WebSocketGroup = "question"
	Like     WebSocketGroup = "like"
)

// this is for the type of server response of the message
type WebSocketType string

const (
	// questions type
	CreateQuestionType      WebSocketType = "createQuestion"
	DeleteQuestionType      WebSocketType = "deleteQuestion"
	EditQuestionType        WebSocketType = "editQuestion"
	AdminDeleteQuestionType WebSocketType = "adminDeleteQuestion"
	AdminEditQuestionType   WebSocketType = "adminEdit"

	// likes type
	ToggleLikeType WebSocketType = "toggleLike"

	// error type
	ErrorType WebSocketType = "error"
)

type WebResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
}

type WebsocketResponse struct {
	Type    WebSocketType  `json:"type,omitempty"`
	Group   WebSocketGroup `json:"group,omitempty"`
	Data    any            `json:"data,omitempty"`
	Error   bool           `json:"error"`
	Message string         `json:"message,omitempty"`
}

func RespondWithError(ctx *gin.Context, code int, errMsg string) {
	err := errors.New(errMsg)

	ctx.Error(err)
	ctx.AbortWithStatusJSON(code, WebResponse{
		Error:   true,
		Message: errMsg,
	})
}

func RespondWithJson(ctx *gin.Context, code int, data any) {
	ctx.JSON(code, WebResponse{
		Error: false,
		Data:  data,
	})
}

func EncodeJson(data any) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("json bind err: ", err)
		return []byte{}
	}
	return jsonData
}

func DecodeJson(data []byte, target any) error {
	err := json.Unmarshal(data, target)
	if err != nil {
		return err
	}
	return nil
}

func WebSocketRespondJson(group WebSocketGroup, t WebSocketType, data any) []byte {
	return EncodeJson(WebsocketResponse{
		Type:  t,
		Group: group,
		Data:  data,
		Error: false,
	})
}

func WebSocketRespondError(group WebSocketGroup, errMsg string) []byte {
	return EncodeJson(WebsocketResponse{
		Type:    ErrorType,
		Group:   group,
		Error:   true,
		Message: errMsg,
	})
}

func CheckNil[t any](anyType t) *t {
	if reflect.ValueOf(anyType).IsZero() {
		return nil
	} else {
		return &anyType
	}
}
