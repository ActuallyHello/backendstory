package enum

import (
	"time"
)

// EnumCreateRequest represents request for creating enum
// @Name EnumCreateRequest
type EnumCreateRequest struct {
	Code  string `json:"code" validate:"required,min=1,max=50"`
	Label string `json:"label" validate:"required,min=1,max=255"`
}

// EnumUpdateRequest represents request for updating enum
// @Name EnumUpdateRequest
type EnumUpdateRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Code  string `json:"code" validate:"omitempty,min=1,max=50"`
	Label string `json:"label" validate:"omitempty,min=1,max=255"`
}

// EnumDTO represents enum data transfer object
// @Name EnumDTO
type EnumDTO struct {
	ID        uint      `json:"id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Code      string    `json:"code" validate:"required"`
	Label     string    `json:"label" validate:"required"`
}

func ToEnumDTO(enum Enum) EnumDTO {
	return EnumDTO{
		ID:    enum.ID,
		Code:  enum.Code,
		Label: enum.Label,
	}
}
