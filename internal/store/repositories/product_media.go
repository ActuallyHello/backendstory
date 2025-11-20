package repositories

import (
	"context"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"gorm.io/gorm"
)

type ProductMediaRepository interface {
	BaseRepository[entities.ProductMedia]

	FindByProductID(ctx context.Context, productMediaID uint) ([]entities.ProductMedia, error)
}

type productMediaRepository struct {
	BaseRepositoryImpl[entities.ProductMedia]
}

func NewProductMediaRepository(db *gorm.DB) *productMediaRepository {
	return &productMediaRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.ProductMedia](db),
	}
}

func (r *productMediaRepository) FindByProductID(ctx context.Context, productID uint) ([]entities.ProductMedia, error) {
	var productMedia []entities.ProductMedia
	if err := r.db.WithContext(ctx).Where("PRODUCTID = ?", productID).Find(&productMedia).Error; err != nil {
		return nil, err
	}
	return productMedia, nil
}
