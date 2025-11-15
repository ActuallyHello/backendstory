package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	BaseRepository[entities.Category]

	FindByCode(ctx context.Context, code string) (entities.Category, error)
	FindByCategoryID(ctx context.Context, categoryID uint) ([]entities.Category, error)
}

type categoryRepository struct {
	BaseRepositoryImpl[entities.Category]
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.Category](db),
	}
}

// FindByCode ищет Enum по коду
func (r *categoryRepository) FindByCode(ctx context.Context, code string) (entities.Category, error) {
	var category entities.Category
	if err := r.db.WithContext(ctx).Where("CODE = ?", code).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Category{}, common.NewNotFoundError("category not found by code")
		}
		return entities.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]entities.Category, error) {
	var categories []entities.Category
	if err := r.db.WithContext(ctx).Where("CATEGORYID = ?", categoryID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
