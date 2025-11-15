package dto

import (
	"time"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

// CategoryCreateRequest represents request for creating category
// @Name CategoryCreateRequest
type CategoryCreateRequest struct {
	Code       string `json:"code" validate:"required,min=1,max=50"`
	Label      string `json:"label" validate:"required,min=1,max=255"`
	CategoryID *uint  `json:"category_id" validate:"omitempty,min=1"`
}

// CategoryUpdateRequest represents request for updating category
// @Name CategoryUpdateRequest
type CategoryUpdateRequest struct {
	ID         uint   `json:"id" validate:"required"`
	Code       string `json:"code" validate:"omitempty,min=1,max=50"`
	Label      string `json:"label" validate:"omitempty,min=1,max=255"`
	CategoryID *uint  `json:"category_id" validate:"omitempty,min=1"`
}

// CategoryDTO represents category data transfer object
// @Name CategoryDTO
type CategoryDTO struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Code       string    `json:"code"`
	Label      string    `json:"label"`
	CategoryID *uint     `json:"category_id"`
}

func ToCategoryDTO(category entities.Category) CategoryDTO {
	var categoryID *uint
	if category.CategoryID.Valid {
		tempCategoryID := uint(category.CategoryID.Int32)
		categoryID = &tempCategoryID
	}

	return CategoryDTO{
		ID:         category.ID,
		CreatedAt:  category.CreatedAt,
		UpdatedAt:  category.UpdatedAt,
		Code:       category.Code,
		Label:      category.Label,
		CategoryID: categoryID,
	}
}
