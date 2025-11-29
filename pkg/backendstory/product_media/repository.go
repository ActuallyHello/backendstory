package productmedia

import (
	"context"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type ProductMediaRepository interface {
	core.BaseRepository[ProductMedia]

	FindByProductID(ctx context.Context, productMediaID uint) ([]ProductMedia, error)
}

type productMediaRepository struct {
	core.BaseRepositoryImpl[ProductMedia]
}

func NewProductMediaRepository(db *gorm.DB) *productMediaRepository {
	return &productMediaRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[ProductMedia](db),
	}
}

func (r *productMediaRepository) FindByProductID(ctx context.Context, productID uint) ([]ProductMedia, error) {
	var productMedia []ProductMedia
	if err := r.GetDB().WithContext(ctx).Where("PRODUCTID = ?", productID).Find(&productMedia).Error; err != nil {
		return nil, err
	}
	return productMedia, nil
}
