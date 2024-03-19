package dtos

import (
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/google/uuid"
)

type QuestionResponse struct {
	QuestionID *uuid.UUID    `json:"question_id,omitempty"`
	EventID    *uuid.UUID    `json:"event_id,omitempty"`
	UserID     *uuid.UUID    `json:"user_id,omitempty"`
	Username   string        `json:"username,omitempty"`
	Content    string        `json:"content,omitempty"`
	Starred    bool          `json:"starred,omitempty"`
	Approved   bool          `json:"approved,omitempty"`
	Answered   bool          `json:"answered,omitempty"`
	LikesCount int           `json:"likes_count"`
	UserLiked  bool          `json:"user_liked"`
	CreatedAt  *time.Time    `json:"created_at,omitempty"`
	UpdatedAt  *time.Time    `json:"updated_at,omitempty"`
	Event      EventResponse `json:"event,omitempty"`
}

type CreateQuestionInput struct {
	EventID  string `json:"event_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Username string `json:"username"`
}

type DeleteQuestionInput struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required"`
}

type EditQuestionInput struct {
	QuestionID string `json:"question_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

func GenerateQuestionResponse(question *models.Question, user User) *QuestionResponse {
	if question == nil {
		return nil
	}

	// Check if the user liked the post
	userLiked := false
	for _, like := range question.Likes {
		if like.UserID == user.ID {
			userLiked = true
			break
		}
	}

	return &QuestionResponse{
		QuestionID: CheckNil(question.QuestionID),
		EventID:    CheckNil(question.EventID),
		UserID:     CheckNil(question.UserID),
		Username:   question.Username,
		Content:    question.Content,
		Starred:    question.Starred,
		Approved:   question.Approved,
		Answered:   question.Answered,
		LikesCount: len(question.Likes),
		UserLiked:  userLiked,
		CreatedAt:  CheckNil(question.CreatedAt),
		UpdatedAt:  CheckNil(question.UpdatedAt),
		Event:      *GenerateEventResponse(&question.Event),
	}
}
