package dtos

import (
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/google/uuid"
)

// because this is a data transfer object
// so everything can be null/empty
// so make everything a pointer except for a type than can be detected by json omitempty as empty value

type MasterResponse struct {
	MasterID  *uuid.UUID `json:"master_id,omitempty"`
	Email     string     `json:"email,omitempty"`
	Type      string     `json:"type,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type MasterSignUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type MasterSignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type MasterOtpInput struct {
	Code string `json:"code" binding:"required"`
}

type InvitationResponse struct {
	InvitationID *uuid.UUID      `json:"invitation_id,omitempty"`
	MasterID     *uuid.UUID      `json:"master_id,omitempty"`
	Code         string          `json:"code,omitempty"`
	ExpireDate   *time.Time      `json:"expire_date,omitempty"`
	Used         bool            `json:"used"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
	Master       *MasterResponse `json:"master,omitempty"`
}

func GenerateMasterResponse(master *models.Master) *MasterResponse {
	if master == nil {
		return nil
	}
	return &MasterResponse{
		MasterID:  CheckNil(master.MasterID),
		Email:     master.Email,
		Type:      "master",
		CreatedAt: CheckNil(master.CreatedAt),
		UpdatedAt: CheckNil(master.UpdatedAt),
	}
}

func GenerateInvitationResponse(invitation *models.Invitation) *InvitationResponse {
	if invitation == nil {
		return nil
	}
	return &InvitationResponse{
		InvitationID: CheckNil(invitation.InvitationID),
		MasterID:     CheckNil(invitation.MasterID),
		Code:         invitation.Code,
		ExpireDate:   CheckNil(invitation.ExpireDate),
		Used:         invitation.Used,
		CreatedAt:    CheckNil(invitation.CreatedAt),
		UpdatedAt:    CheckNil(invitation.Master.UpdatedAt),
		Master:       GenerateMasterResponse(&invitation.Master),
	}
}
