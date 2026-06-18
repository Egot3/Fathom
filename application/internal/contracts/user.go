package contracts

import (
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type DeleteUserRequest struct {
	UserUUID uuid.UUID `json:"uuid"`
}

type GetUserRequest struct {
	UserUUID uuid.UUID `json:"uuid"`
}

type GetUserResponse struct {
	User models.User `json:"user"`
}

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User models.User `json:"user"`
}

type RegisterRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User models.User `json:"user"`
}

type PatchRequest struct {
	UUID      uuid.UUID `json:"uuid"`
	Nickname  *string   `json:"nickname,omitempty"`
	Password  *string   `json:"password,omitempty"`
	IsTeacher *bool     `json:"is_teacher,omitempty"`
}

type ListUsersRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ListUsersResponse struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`

	Users []models.User `json:"users"`
}
