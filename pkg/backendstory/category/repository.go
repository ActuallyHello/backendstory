package category

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	core.BaseRepository[Category]

	FindByCode(ctx context.Context, code string) (Category, error)
	FindByCategoryID(ctx context.Context, categoryID uint) ([]Category, error)
}

type categoryRepository struct {
	core.BaseRepositoryImpl[Category]
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[Category](db),
	}
}

// FindByCode ищет Enum по коду
func (r *categoryRepository) FindByCode(ctx context.Context, code string) (Category, error) {
	var category Category
	if err := r.GetDB().WithContext(ctx).Where("CODE = ?", code).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Category{}, core.NewNotFoundError("category not found by code")
		}
		return Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]Category, error) {
	var categories []Category
	if err := r.GetDB().WithContext(ctx).Where("CATEGORYID = ?", categoryID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
