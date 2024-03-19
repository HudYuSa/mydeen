package dtos

import (
	"time"

	"github.com/HudYuSa/mydeen/db/models"
	"github.com/google/uuid"
)

type AdminResponse struct {
	AdminID      *uuid.UUID `json:"admin_id,omitempty"`
	InvitationID *uuid.UUID `json:"invitation_id,omitempty"`
	Username     string     `json:"username,omitempty"`
	Email        string     `json:"email,omitempty"`
	AdminCode    string     `json:"admin_code,omitempty"`
	Type         string     `json:"type,omitempty"`
	Enable2fa    bool       `json:"enable2fa"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type AdminSignUpInput struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	InvitationCode string `json:"invitation_code" binding:"required"`
}

type AdminSignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminOtpInput struct {
	Code string `json:"code" binding:"required"`
}

type UpdateAdminUsername struct {
	Username string `json:"username" binding:"required"`
}

type UpdateAdminEmail struct {
	Email string `json:"email" binding:"required,email"`
}

func GenerateAdminResponse(admin *models.Admin) *AdminResponse {
	if admin == nil {
		return nil
	}
	return &AdminResponse{
		AdminID:      CheckNil(admin.AdminID),
		InvitationID: CheckNil(admin.InvitationID),
		Username:     admin.Username,
		Email:        admin.Email,
		AdminCode:    admin.AdminCode,
		Type:         "admin",
		Enable2fa:    admin.Enable2fa,
		CreatedAt:    CheckNil(admin.CreatedAt),
		UpdatedAt:    CheckNil(admin.UpdatedAt),
	}
}
