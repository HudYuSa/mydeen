package dtos

import (
	"github.com/HudYuSa/mydeen/db/models"
	"github.com/google/uuid"
)

type LikeResponse struct {
	LikeID     *uuid.UUID `json:"like_id,omitempty"`
	QuestionID *uuid.UUID `json:"question_id,omitempty"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Liked      bool       `json:"liked"`
}

type ToggleLikeInput struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required"`
}

func GenerateLikeResponse(like *models.Like, liked bool) *LikeResponse {
	if like == nil {
		return nil
	}

	return &LikeResponse{
		LikeID:     CheckNil(like.LikeID),
		QuestionID: CheckNil(like.QuestionID),
		UserID:     CheckNil(like.UserID),
		Liked:      liked,
	}
}
