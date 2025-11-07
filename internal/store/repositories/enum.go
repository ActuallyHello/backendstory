package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type EnumRepository interface {
	BaseRepository[entities.Enum]

	FindByCode(ctx context.Context, code string) (entities.Enum, error)
}

type enumRepository struct {
	BaseRepositoryImpl[entities.Enum]
}

func NewEnumRepository(db *gorm.DB) *enumRepository {
	return &enumRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.Enum](db),
	}
}

// FindByCode ищет Enum по коду
func (r *enumRepository) FindByCode(ctx context.Context, code string) (entities.Enum, error) {
	var enum entities.Enum
	if err := r.db.WithContext(ctx).Where("CODE = ?", code).First(&enum).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Enum{}, common.NewNotFoundError("enum not found by code")
		}
		return entities.Enum{}, err
	}
	return enum, nil
}
