package dtos

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

// this package is where u put all data transfer object representation
// and all its function

type WebResponse struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error"`
	Data    any    `json:"data,omitempty"`
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

func CheckNil[t any](anyType t) *t {
	if reflect.ValueOf(anyType).IsZero() {
		return nil
	} else {
		return &anyType
	}
}
