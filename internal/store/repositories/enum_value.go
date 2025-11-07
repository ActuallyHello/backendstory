package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type EnumValueRepository interface {
	BaseRepository[entities.EnumValue]

	FindByEnumID(ctx context.Context, enumerationID uint) ([]entities.EnumValue, error)
	FindByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (entities.EnumValue, error)
}

type enumValueRepository struct {
	BaseRepositoryImpl[entities.EnumValue]
}

func NewEnumValueRepository(db *gorm.DB) *enumValueRepository {
	return &enumValueRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.EnumValue](db),
	}
}

// FindByCode ищет EnumValue по коду
func (r *enumValueRepository) FindByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (entities.EnumValue, error) {
	var enumValue entities.EnumValue
	if err := r.db.WithContext(ctx).Where("CODE = ?", code).Where("ENUMERATIONID = ?", enumerationID).First(&enumValue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.EnumValue{}, common.NewNotFoundError("enum value not found by code")
		}
		return entities.EnumValue{}, err
	}
	return enumValue, nil
}

// FindByEnumerationID ищет все EnumValue по EnumerationID
func (r *enumValueRepository) FindByEnumID(ctx context.Context, enumerationID uint) ([]entities.EnumValue, error) {
	var enumValues []entities.EnumValue
	if err := r.db.WithContext(ctx).Where("ENUMERATIONID = ?", enumerationID).Find(&enumValues).Error; err != nil {
		return nil, err
	}
	return enumValues, nil
}
