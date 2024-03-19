package routes

import (
	"github.com/HudYuSa/mydeen/pkg/controllers"
	"github.com/HudYuSa/mydeen/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type EventRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type eventRoutes struct {
	EventController controllers.EventController
}

func NewEventRoutes(eventController controllers.EventController) EventRoutes {
	return &eventRoutes{
		EventController: eventController,
	}
}

func (er *eventRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/event")

	router.GET(("/search/:event_code"), er.EventController.GetEvent)
	router.GET("/live/:admin_code", er.EventController.GetLiveEvent)

	router.Use(middlewares.AuthenticateAdmin())
	router.POST("", er.EventController.CreateEvent)
	router.GET(("/all"), er.EventController.GetAdminEvents)
	router.GET(("/scheduled"), er.EventController.GetScheduledAdminEvents)
	router.GET(("/finished"), er.EventController.GetFinishedAdminEvents)
	router.DELETE("/:event_id", er.EventController.DeleteEvent)
	router.GET(("/:event_id/start"), er.EventController.StartEvent)
	router.GET("/:event_id/finish", er.EventController.FinishEvent)
	router.PATCH("/:event_id/event-name", er.EventController.UpdateName)
	router.PATCH("/:event_id/start-date", er.EventController.UpdateDate)
	router.PATCH("/:event_id/moderation", er.EventController.UpdateModeration)
	router.PATCH("/:event_id/max-question-length", er.EventController.UpdateMaxQuestionLength)
	router.PATCH("/:event_id/max-questions", er.EventController.UpdateMaxQuestions)
}
