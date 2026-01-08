package enumvalue

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type EnumValueRepository interface {
	core.BaseRepository[EnumValue]

	FindByEnumID(ctx context.Context, enumerationID uint) ([]EnumValue, error)
	FindByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (EnumValue, error)
}

type enumValueRepository struct {
	core.BaseRepositoryImpl[EnumValue]
}

func NewEnumValueRepository(db *gorm.DB) *enumValueRepository {
	return &enumValueRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[EnumValue](db),
	}
}

// FindByCode ищет EnumValue по коду
func (r *enumValueRepository) FindByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (EnumValue, error) {
	var enumValue EnumValue
	if err := r.GetDB(ctx).Where("CODE = ?", code).Where("ENUMERATIONID = ?", enumerationID).First(&enumValue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EnumValue{}, core.NewNotFoundError("Значение перечисления не существует")
		}
		return EnumValue{}, err
	}
	return enumValue, nil
}

// FindByEnumerationID ищет все EnumValue по EnumerationID
func (r *enumValueRepository) FindByEnumID(ctx context.Context, enumerationID uint) ([]EnumValue, error) {
	var enumValues []EnumValue
	if err := r.GetDB(ctx).Where("ENUMERATIONID = ?", enumerationID).Find(&enumValues).Error; err != nil {
		return nil, err
	}
	return enumValues, nil
}
