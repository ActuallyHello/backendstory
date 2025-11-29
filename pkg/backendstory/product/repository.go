package product

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type ProductRepository interface {
	core.BaseRepository[Product]

	FindByCode(ctx context.Context, code string) (Product, error)
	FindBySku(ctx context.Context, sku string) (Product, error)
	FindByCategoryID(ctx context.Context, categoryID uint) ([]Product, error)
}

type productRepository struct {
	core.BaseRepositoryImpl[Product]
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[Product](db),
	}
}

// FindByCode ищет по коду
func (r *productRepository) FindByCode(ctx context.Context, code string) (Product, error) {
	var product Product
	if err := r.GetDB().WithContext(ctx).Where("CODE = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Product{}, core.NewNotFoundError("product not found by code")
		}
		return Product{}, err
	}
	return product, nil
}

// FindBySku ищет по артиклу
func (r *productRepository) FindBySku(ctx context.Context, sku string) (Product, error) {
	var product Product
	if err := r.GetDB().WithContext(ctx).Where("SKU = ?", sku).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Product{}, core.NewNotFoundError("product not found by sku")
		}
		return Product{}, err
	}
	return product, nil
}

// FindByCategoryID ищеn по категории
func (r *productRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]Product, error) {
	var products []Product
	if err := r.GetDB().WithContext(ctx).Where("CATEGORYID = ?", categoryID).Find(&products).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("product not found by category id")
		}
		return nil, err
	}
	return products, nil
}
