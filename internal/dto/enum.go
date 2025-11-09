package dto

import (
	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

type EnumCreateRequest struct {
	Code  string `json:"code" validate:"required,min=1,max=50"`
	Label string `json:"label" validate:"required,min=1,max=255"`
}

type EnumUpdateRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Code  string `json:"code" validate:"omitempty,min=1,max=50"`
	Label string `json:"label" validate:"omitempty,min=1,max=255"`
}

type EnumDTO struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

func ToEnumDTO(enum entities.Enum) EnumDTO {
	return EnumDTO{
		ID:    enum.ID,
		Code:  enum.Code,
		Label: enum.Label,
	}
}
