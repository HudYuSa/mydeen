package dtos

import (
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/google/uuid"
)

type EventResponse struct {
	EventID           *uuid.UUID            `json:"event_id,omitempty"`
	AdminId           *uuid.UUID            `json:"admin_id,omitempty"`
	EventName         string                `json:"event_name,omitempty"`
	Status            models.Status         `json:"status,omitempty"`
	Moderation        bool                  `json:"moderation"`
	MaxQuestions      models.MaxQuestions   `json:"max_questions,omitempty"`
	MaxQuestionLength models.QuestionLength `json:"max_question_length,omitempty"`
	EventCode         string                `json:"event_code,omitempty"`
	StartDate         *time.Time            `json:"start_date,omitempty"`
	CreatedAt         *time.Time            `json:"created_at,omitempty"`
	UpdatedAt         *time.Time            `json:"updated_at,omitempty"`
	Admin             *AdminResponse        `json:"admin,omitempty"`
}

type CreateEventInput struct {
	EventName string `json:"event_name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
}

type UpdateEventNameInput struct {
	EventName string `json:"event_name" binding:"required"`
}

type UpdateEventDateInput struct {
	StartDate string `json:"start_date" binding:"required"`
}

type UpdateModerationInput struct {
	Moderation bool `json:"moderation"`
}

type UpdateMaxQuestionLengthInput struct {
	MaxQuestionLength int `json:"max_question_length" binding:"required"`
}

type UpdateMaxQuestions struct {
	MaxQuestions int `json:"max_questions" binding:"required"`
}

func GenerateEventResponse(event *models.Event) *EventResponse {
	if event == nil {
		return nil
	}
	return &EventResponse{
		EventID:           CheckNil(event.EventID),
		AdminId:           CheckNil(event.AdminID),
		EventName:         event.EventName,
		Status:            event.Status,
		Moderation:        event.Moderation,
		MaxQuestions:      event.MaxQuestions,
		MaxQuestionLength: event.MaxQuestionLength,
		EventCode:         event.EventCode,
		StartDate:         CheckNil(event.StartDate),
		CreatedAt:         CheckNil(event.CreatedAt),
		UpdatedAt:         CheckNil(event.UpdatedAt),
		Admin:             GenerateAdminResponse(&event.Admin),
	}
}
