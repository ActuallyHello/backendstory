package enumvalue

import (
	"time"
)

// EnumValueCreateRequest represents request for creating enum value
// @Name EnumValueCreateRequest
type EnumValueCreateRequest struct {
	Code   string `json:"code" validate:"required,min=1,max=50"`
	Label  string `json:"label" validate:"required,min=1,max=255"`
	EnumID uint   `json:"enumeration_id" validate:"required,gt=0"`
}

// EnumValueUpdateRequest represents request for updating enum value
// @Name EnumValueUpdateRequest
type EnumValueUpdateRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Code  string `json:"code" validate:"omitempty,min=1,max=50"`
	Label string `json:"label" validate:"omitempty,min=1,max=255"`
}

// EnumValueDTO represents enum value data transfer object
// @Name EnumValueDTO
type EnumValueDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Code      string    `json:"code"`
	Label     string    `json:"label"`
	EnumID    uint      `json:"enum_id"`
}

func ToEnumValueDTO(enumValue EnumValue) EnumValueDTO {
	return EnumValueDTO{
		ID:        enumValue.ID,
		CreatedAt: enumValue.CreatedAt,
		UpdatedAt: enumValue.UpdatedAt,
		Code:      enumValue.Code,
		Label:     enumValue.Label,
		EnumID:    enumValue.EnumID,
	}
}
