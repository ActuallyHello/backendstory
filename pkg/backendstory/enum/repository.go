package enum

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type EnumRepository interface {
	core.BaseRepository[Enum]

	FindByCode(ctx context.Context, code string) (Enum, error)
}

type enumRepository struct {
	core.BaseRepositoryImpl[Enum]
}

func NewEnumRepository(db *gorm.DB) *enumRepository {
	return &enumRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[Enum](db),
	}
}

func (r *enumRepository) FindByCode(ctx context.Context, code string) (Enum, error) {
	var enum Enum
	if err := r.GetDB(ctx).Where("CODE = ?", code).First(&enum).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Enum{}, core.NewNotFoundError("Перечисление не найдено по переданному коду")
		}
		return Enum{}, err
	}
	return enum, nil
}
