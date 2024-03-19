package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/HudYuSa/mydeen/pkg/dtos"
	"github.com/HudYuSa/mydeen/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventController interface {
	CreateEvent(ctx *gin.Context)
	GetEvent(ctx *gin.Context)
	GetAdminEvents(ctx *gin.Context)
	GetAdminEvent(ctx *gin.Context)
	GetScheduledAdminEvents(ctx *gin.Context)
	GetFinishedAdminEvents(ctx *gin.Context)
	GetLiveEvent(ctx *gin.Context)
	UpdateEvent(ctx *gin.Context)
	DeleteEvent(ctx *gin.Context)
	StartEvent(ctx *gin.Context)
	FinishEvent(ctx *gin.Context)
	UpdateName(ctx *gin.Context)
	UpdateDate(ctx *gin.Context)
	UpdateModeration(ctx *gin.Context)
	UpdateMaxQuestionLength(ctx *gin.Context)
	UpdateMaxQuestions(ctx *gin.Context)
}

type eventController struct {
	DB *gorm.DB
}

func NewEventController(db *gorm.DB) EventController {
	return &eventController{
		DB: db,
	}
}

func (ec *eventController) CreateEvent(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	var payload dtos.CreateEventInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// fmt.Println(payload)
	// create event entity
	now := time.Now().UTC()
	date, err := time.Parse("2006-01-02 15:04:05", payload.StartDate)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	randomCode, err := utils.GenerateRandomNumCode()
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, err.Error())
	}

	fmt.Println(randomCode)

	newEvent := models.Event{
		AdminID:           currentAdmin.AdminID,
		EventName:         payload.EventName,
		Status:            models.Scheluded,
		Moderation:        false,
		MaxQuestions:      models.MidCount,
		MaxQuestionLength: models.VeryLong,
		EventCode:         randomCode,
		StartDate:         date,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// save to database
	eventResult := ec.DB.WithContext(dbTimeoutCtx).Create(&newEvent)
	if eventResult.Error != nil && strings.Contains(eventResult.Error.Error(), "duplicate key value violates unique") {
		dtos.RespondWithError(ctx, http.StatusConflict, "Duplicate event code")
		return
	} else if eventResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		return
	}

	// send the response
	dtos.RespondWithJson(ctx, http.StatusCreated, dtos.GenerateEventResponse(&newEvent))
}

func (ec *eventController) GetEvent(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	eventCode := ctx.Param("event_code")

	event := models.Event{}
	eventResult := ec.DB.WithContext(dbTimeoutCtx).Where("event_code = ?", eventCode).First(&event)
	if eventResult.Error != nil {
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given code")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	dtos.RespondWithJson(ctx, http.StatusOK, dtos.GenerateEventResponse(&event))
}

func (ec *eventController) GetAdminEvents(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	// get events based on admin id
	events := []models.Event{}
	eventsResult := ec.DB.WithContext(dbTimeoutCtx).Where("admin_id = ?", currentAdmin.AdminID).Order("created_at DESC").Find(&events)
	if eventsResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, eventsResult.Error.Error())
		return
	}

	eventsResponse := []dtos.EventResponse{}

	for _, event := range events {
		eventsResponse = append(eventsResponse, *dtos.GenerateEventResponse(&event))
	}

	dtos.RespondWithJson(ctx, http.StatusOK, eventsResponse)
}

func (ec *eventController) GetAdminEvent(ctx *gin.Context) {
	// create a function to get crucial admin data but only the admin that have it is allowed to access
	panic("not implemented") // TODO: Implement
}

func (ec *eventController) GetScheduledAdminEvents(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	// get events based on admin id
	events := []models.Event{}
	eventsResult := ec.DB.WithContext(dbTimeoutCtx).Where("admin_id = ? AND status = ?", currentAdmin.AdminID, models.Scheluded).Order("created_at DESC").Find(&events)
	if eventsResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, eventsResult.Error.Error())
		return
	}

	eventsResponse := []dtos.EventResponse{}

	for _, event := range events {
		eventsResponse = append(eventsResponse, *dtos.GenerateEventResponse(&event))
	}

	fmt.Println(eventsResponse)
	dtos.RespondWithJson(ctx, http.StatusOK, eventsResponse)
}

func (ec *eventController) GetFinishedAdminEvents(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	// get events based on admin id
	events := []models.Event{}
	eventsResult := ec.DB.WithContext(dbTimeoutCtx).Where("admin_id = ? AND status = ?", currentAdmin.AdminID, models.Finished).Order("created_at DESC").Find(&events)
	if eventsResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, eventsResult.Error.Error())
		return
	}

	eventsResponse := []dtos.EventResponse{}

	for _, event := range events {
		eventsResponse = append(eventsResponse, *dtos.GenerateEventResponse(&event))
	}
	fmt.Println(eventsResponse)

	dtos.RespondWithJson(ctx, http.StatusOK, eventsResponse)
}

func (ec *eventController) GetLiveEvent(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)

	adminCode := ctx.Param("admin_code")

	// cari admin berdasarkan admin codenya di database
	admin := models.Admin{}
	adminResult := ec.DB.WithContext(dbTimeoutCtx).Where("admin_code = ?", adminCode).First(&admin)
	if adminResult.Error != nil {
		switch adminResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no admin with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, adminResult.Error.Error())
		}
		return
	}

	// cari event dari adminnya yang sedang live
	events := []models.Event{}
	eventResult := ec.DB.WithContext(dbTimeoutCtx).Where("admin_id = ? AND status = ?", admin.AdminID, models.Live).Find(&events)
	if eventResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		return
	}

	eventsResponse := []dtos.EventResponse{}

	fmt.Println(eventsResponse)

	for _, event := range events {
		eventsResponse = append(eventsResponse, *dtos.GenerateEventResponse(&event))
	}

	// kirim response
	dtos.RespondWithJson(ctx, http.StatusOK, eventsResponse)
}

func (ec *eventController) UpdateEvent(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ec *eventController) DeleteEvent(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")

	// get event by event_id
	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	deleteEventResult := ec.DB.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).Delete(&models.Event{})
	if deleteEventResult.Error != nil {
		dtos.RespondWithError(ctx, http.StatusInternalServerError, deleteEventResult.Error.Error())
		return
	}

	fmt.Println(deleteEventResult.RowsAffected)
	if deleteEventResult.RowsAffected < 1 {
		dtos.RespondWithError(ctx, http.StatusBadRequest, "there is no event with the given id")
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully delete event")
}

func (ec *eventController) StartEvent(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")

	// check for another live events
	checkEvents := []models.Event{}
	checkEventsResult := tx.WithContext(dbTimeoutCtx).Where("status = ? AND admin_id = ?", models.Live, currentAdmin.AdminID).Find(&checkEvents)
	if checkEventsResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, checkEventsResult.Error.Error())
		return
	}

	// check if there's another live event
	if checkEventsResult.RowsAffected > 0 {
		tx.Rollback()
		fmt.Println(checkEventsResult.RowsAffected, "live events")
		dtos.RespondWithError(ctx, http.StatusBadRequest, "There is an ongoing live event")
		return
	}

	// get event by event_id
	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	// check if event status
	if event.Status == models.Live {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusBadRequest, "Your event has already started")
		return
	} else if event.Status == models.Finished {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusBadRequest, "Your event already finished")
		return
	}

	// update the status
	event.Status = models.Live

	UpdateEventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", event.EventID).Save(&event)
	if UpdateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully started your event")
}

func (ec *eventController) FinishEvent(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")

	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	fmt.Println(event.Status)
	// check if event status
	if event.Status == models.Finished {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusBadRequest, "Your event already finished")
		return
	} else if event.Status == models.Scheluded {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusBadRequest, "Your event hasn't started yet")
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	// update the status
	event.Status = models.Finished

	UpdateEventResult := ec.DB.WithContext(dbTimeoutCtx).Where("event_id = ?", event.EventID).Save(&event)
	// UpdateEventResult := ec.DB.WithContext(dbTimeoutCtx).Model(&models.Event{}).Where("event_id = ? AND status = ?", eventId, models.Live).Update("status", models.Finished)
	if UpdateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully finished your event")
}

func (ec *eventController) UpdateName(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")
	var payload dtos.UpdateEventNameInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(payload)

	// get event by event_id
	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	// update the name
	event.EventName = payload.EventName

	updateEventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).Save(&event)
	if updateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, updateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully update event name")
}

func (ec *eventController) UpdateDate(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")

	// try to bind the request body to the payload struct
	var payload dtos.UpdateEventDateInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// get event by event_id
	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	date, err := time.Parse("2006-01-02 15:04:05", payload.StartDate)
	if err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// update event date
	event.StartDate = date

	UpdateEventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).Save(&event)
	if UpdateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully update event date")
}

func (ec *eventController) UpdateModeration(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")
	var payload dtos.UpdateModerationInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// get event data
	event := models.Event{}
	checkEventsResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if checkEventsResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, checkEventsResult.Error.Error())
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	// update event moderation
	event.Moderation = payload.Moderation

	UpdateEventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", event.EventID).Save(&event)
	if UpdateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully update event moderation")
}

func (ec *eventController) UpdateMaxQuestionLength(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")
	var payload dtos.UpdateMaxQuestionLengthInput

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// get event by event_id
	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	// update question max length
	event.MaxQuestionLength = models.QuestionLength(payload.MaxQuestionLength)

	UpdateEventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", event.EventID).Save(&event)
	if UpdateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully update event max question length")
}

func (ec *eventController) UpdateMaxQuestions(ctx *gin.Context) {
	dbTimeoutCtx := ctx.MustGet("dbTimeoutContext").(context.Context)
	currentAdmin := ctx.MustGet("currentAdmin").(models.Admin)

	tx := ec.DB.Begin()

	eventId := ctx.Param("event_id")
	var payload dtos.UpdateMaxQuestions

	// try to bind the request body to the payload struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		dtos.RespondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// get event by event_id
	event := models.Event{}
	eventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).First(&event)
	if eventResult.Error != nil {
		tx.Rollback()
		switch eventResult.Error.Error() {
		case "record not found":
			dtos.RespondWithError(ctx, http.StatusNotFound, "there is no event with the given id")
		default:
			dtos.RespondWithError(ctx, http.StatusInternalServerError, eventResult.Error.Error())
		}
		return
	}

	// check if admin is the admin that created the event
	if event.AdminID != currentAdmin.AdminID {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusUnauthorized, "You're not allowed to access this endpoint")
		return
	}

	// update max question ask
	event.MaxQuestions = models.MaxQuestions(payload.MaxQuestions)

	UpdateEventResult := tx.WithContext(dbTimeoutCtx).Where("event_id = ?", eventId).Save(&event)
	if UpdateEventResult.Error != nil {
		tx.Rollback()
		dtos.RespondWithError(ctx, http.StatusInternalServerError, UpdateEventResult.Error.Error())
		return
	}

	tx.Commit()

	dtos.RespondWithJson(ctx, http.StatusOK, "Successfully update event max questions")
}
