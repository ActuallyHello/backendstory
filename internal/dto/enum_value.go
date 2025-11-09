package dto

import (
	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

type EnumValueCreateRequest struct {
	Code   string `json:"code" validate:"required,min=1,max=50"`
	Label  string `json:"label" validate:"required,min=1,max=255"`
	EnumID uint   `json:"enumeration_id" validate:"required,gt=0"`
}

type EnumValueUpdateRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Code  string `json:"code" validate:"omitempty,min=1,max=50"`
	Label string `json:"label" validate:"omitempty,min=1,max=255"`
}

type EnumValueDTO struct {
	ID     uint   `json:"id"`
	Code   string `json:"code"`
	Label  string `json:"label"`
	EnumID uint   `json:"enum_id"`
}

func ToEnumValueDTO(enumValue entities.EnumValue) EnumValueDTO {
	return EnumValueDTO{
		ID:     enumValue.ID,
		Code:   enumValue.Code,
		Label:  enumValue.Label,
		EnumID: enumValue.EnumID,
	}
}
