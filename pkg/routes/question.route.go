package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type QuestionRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type questionRoutes struct {
	QuestionController controllers.QuestionController
}

func NewQuestionRoutes(questionController controllers.QuestionController) QuestionRoutes {
	return &questionRoutes{
		QuestionController: questionController,
	}
}

func (qr *questionRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/questions")

	router.GET("/:event_id", qr.QuestionController.GetEventQuestions)
}
