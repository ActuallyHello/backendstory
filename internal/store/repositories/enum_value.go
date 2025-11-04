package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type EnumValueRepository interface {
	Create(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error)
	Update(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error)
	Delete(ctx context.Context, enumValue entities.EnumValue) error

	FindAll(ctx context.Context) ([]entities.EnumValue, error)
	FindById(ctx context.Context, id uint) (entities.EnumValue, error)
	FindByEnumID(ctx context.Context, enumerationID uint) ([]entities.EnumValue, error)
	FindByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (entities.EnumValue, error)
}

type enumValueRepository struct {
	db *gorm.DB
}

func NewEnumValueRepository(db *gorm.DB) *enumValueRepository {
	return &enumValueRepository{db: db}
}

// Create создает новую запись EnumValue
func (r *enumValueRepository) Create(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error) {
	if err := r.db.WithContext(ctx).Create(&enumValue).Error; err != nil {
		return entities.EnumValue{}, err
	}
	return enumValue, nil
}

// Delete выполняет удаление EnumValue
func (r *enumValueRepository) Delete(ctx context.Context, enumValue entities.EnumValue) error {
	if err := r.db.WithContext(ctx).Delete(&enumValue).Error; err != nil {
		return err
	}
	return nil
}

// Update обновляет существующую запись EnumValue
func (r *enumValueRepository) Update(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error) {
	if err := r.db.WithContext(ctx).Save(&enumValue).Error; err != nil {
		return entities.EnumValue{}, err
	}
	return enumValue, nil
}

// FindAll ищет все EnumValue
func (r *enumValueRepository) FindAll(ctx context.Context) ([]entities.EnumValue, error) {
	var enumValues []entities.EnumValue
	if err := r.db.WithContext(ctx).Find(&enumValues).Error; err != nil {
		return nil, err
	}
	return enumValues, nil
}

// FindById ищет EnumValue по ID
func (r *enumValueRepository) FindById(ctx context.Context, id uint) (entities.EnumValue, error) {
	var enumValue entities.EnumValue
	if err := r.db.WithContext(ctx).First(&enumValue, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.EnumValue{}, common.NewNotFoundError("enum value not found by id")
		}
		return entities.EnumValue{}, err
	}
	return enumValue, nil
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
