package dtos

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID
}

type UserResponse struct {
	ID *uuid.UUID `json:"id"`
}

func GenerateUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID: CheckNil(user.ID),
	}
}
