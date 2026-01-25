package product

import (
	"time"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
)

// ProductCreateRequest represents request for creating product
// @Name ProductCreateRequest
type ProductCreateRequest struct {
	Code       string `json:"code" validate:"required,min=1,max=50"`
	Label      string `json:"label" validate:"required,min=1,max=255"`
	Sku        string `json:"sku" validate:"required,min=1,max=255"`
	Price      string `json:"price" validate:"required,min=1,max=30"`
	Quantity   uint   `json:"quantity" validate:"required,gte=0"`
	IsVisible  bool   `json:"is_visible" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required,min=1"`
}

type ProductStatusChangeRequest struct {
	ID         uint   `json:"id" validate:"required"`
	StatusCode string `json:"status_code" validate:"required,min=1,max=255"`
}

type ProductPriceChangeRequest struct {
	ID    uint   `json:"id" validate:"required"`
	Price string `json:"price" validate:"required,min=1,max=30"`
}

// ProductUpdateRequest represents request for updating product
// @Name ProductUpdateRequest
type ProductUpdateRequest struct {
	ID         uint   `json:"id" validate:"required"`
	Code       string `json:"code" validate:"omitempty,min=1,max=50"`
	Label      string `json:"label" validate:"omitempty,min=1,max=255"`
	Sku        string `json:"sku" validate:"required,min=1,max=255"`
	Price      string `json:"price" validate:"required,min=1,max=30"`
	Quantity   uint   `json:"quantity" validate:"required,gte=0"`
	CategoryID uint   `json:"category_id" validate:"omitempty,min=1,max=255"`
	StatusID   uint   `json:"status_id" validate:"required,gt=0"`
}

// ProductDTO represents product data transfer object
// @Name ProductDTO
type ProductDTO struct {
	ID         uint                   `json:"id"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	DeletedAt  *time.Time             `json:"delted_at"`
	Code       string                 `json:"code"`
	Label      string                 `json:"label"`
	Sku        string                 `json:"sku"`
	Price      string                 `json:"price"`
	Quantity   uint                   `json:"quantity"`
	CategoryID uint                   `json:"category_id"`
	StatusDTO  enumvalue.EnumValueDTO `json:"status_dto"`
	IsVisible  bool                   `json:"is_visible"`
}

func ToProductDTO(product Product, status enumvalue.EnumValueDTO) ProductDTO {
	var deletedAt *time.Time
	if product.DeletedAt.Valid {
		deletedAt = &product.DeletedAt.Time
	}
	return ProductDTO{
		ID:         product.ID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
		DeletedAt:  deletedAt,
		Sku:        product.Sku,
		Price:      product.Price.String(),
		Quantity:   product.Quantity,
		Code:       product.Code,
		Label:      product.Label,
		CategoryID: product.CategoryID,
		StatusDTO:  status,
		IsVisible:  product.IsVisible,
	}
}
