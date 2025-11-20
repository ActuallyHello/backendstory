package dto

import (
	"time"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

type ProductMediaCreateRequest struct {
	ProductID uint `json:"productId" validate:"required,min=1"`
}

type ProductMediaDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Link      string    `json:"link"`
	ProductID uint      `json:"product_id"`
}

func ToProductMediaDTO(media entities.ProductMedia) ProductMediaDTO {
	return ProductMediaDTO{
		ID:        media.ID,
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
		Link:      media.Link,
		ProductID: media.ProductID,
	}
}
