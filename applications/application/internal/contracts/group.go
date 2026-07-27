package contracts

import (
	"github.com/egot3/fathom/internal/models"
	"github.com/google/uuid"
)

type PostGroupRequest struct {
	Name       string      `json:"name"`
	Appendants []uuid.UUID `json:"appendants"`
}

type GetGroupResponse struct {
	Group models.Group `json:"group"`
}

type PatchGroupRequest struct {
	Name *string `json:"name"`
}

type AppendUsersRequest struct {
	Appendants []uuid.UUID `json:"appendants"`
}

type RemoveUsersRequest struct {
	Removants []uuid.UUID `json:"removants"`
}

type ListGroupsResponse struct {
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Total  int            `json:"total"`
	Groups []models.Group `json:"groups"`
}
