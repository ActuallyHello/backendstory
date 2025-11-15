package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[entities.Product]

	FindByCode(ctx context.Context, code string) (entities.Product, error)
	FindBySku(ctx context.Context, sku string) (entities.Product, error)
	FindByCategoryID(ctx context.Context, categoryID uint) ([]entities.Product, error)
}

type productRepository struct {
	BaseRepositoryImpl[entities.Product]
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.Product](db),
	}
}

// FindByCode ищет по коду
func (r *productRepository) FindByCode(ctx context.Context, code string) (entities.Product, error) {
	var product entities.Product
	if err := r.db.WithContext(ctx).Where("CODE = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Product{}, common.NewNotFoundError("product not found by code")
		}
		return entities.Product{}, err
	}
	return product, nil
}

// FindBySku ищет по артиклу
func (r *productRepository) FindBySku(ctx context.Context, sku string) (entities.Product, error) {
	var product entities.Product
	if err := r.db.WithContext(ctx).Where("SKU = ?", sku).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Product{}, common.NewNotFoundError("product not found by sku")
		}
		return entities.Product{}, err
	}
	return product, nil
}

// FindByCategoryID ищеn по категории
func (r *productRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]entities.Product, error) {
	var products []entities.Product
	if err := r.db.WithContext(ctx).Where("CATEGORYID = ?", categoryID).Find(&products).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("product not found by category id")
		}
		return nil, err
	}
	return products, nil
}
