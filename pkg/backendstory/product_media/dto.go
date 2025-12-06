package productmedia

import (
	"time"
)

// ProductMediaCreateRequest represents request for creating product media
// @Description Request body for creating product media
type ProductMediaCreateRequest struct {
	ProductID uint `json:"productId" validate:"required,min=1" example:"123"`
}

// ProductMediaDTO represents product media data transfer object
// @Description Product media information
type ProductMediaDTO struct {
	ID        uint      `json:"id" example:"1" validate:"required"`
	CreatedAt time.Time `json:"created_at" example:"2023-10-05T14:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-10-05T14:30:00Z"`
	Link      string    `json:"link" example:"/static/media/product_1.jpg"`
	ProductID uint      `json:"product_id" example:"123" validate:"required"`
}

func ToProductMediaDTO(media ProductMedia) ProductMediaDTO {
	return ProductMediaDTO{
		ID:        media.ID,
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
		Link:      media.Link,
		ProductID: media.ProductID,
	}
}
