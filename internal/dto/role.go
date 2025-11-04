package dto

import "time"

type RoleCreateRequest struct {
	Code  string `json:"code" validate:"required,min=1,max=50"`
	Label string `json:"label" validate:"required,min=1,max=255"`
}

type RoleUpdateRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Code  string `json:"code" validate:"omitempty,min=1,max=50"`
	Label string `json:"label" validate:"omitempty,min=1,max=255"`
}

type RoleDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Code      string    `json:"code"`
	Label     string    `json:"label"`
}
