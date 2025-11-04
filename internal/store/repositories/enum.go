package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type EnumRepository interface {
	Create(ctx context.Context, enum entities.Enum) (entities.Enum, error)
	Update(ctx context.Context, enum entities.Enum) (entities.Enum, error)
	Delete(ctx context.Context, enum entities.Enum) error

	FindAll(ctx context.Context) ([]entities.Enum, error)
	FindById(ctx context.Context, id uint) (entities.Enum, error)
	FindByCode(ctx context.Context, code string) (entities.Enum, error)
}

type enumRepository struct {
	db *gorm.DB
}

func NewEnumRepository(db *gorm.DB) *enumRepository {
	return &enumRepository{db: db}
}

// Create создает новую запись Enum
func (r *enumRepository) Create(ctx context.Context, enum entities.Enum) (entities.Enum, error) {
	if err := r.db.WithContext(ctx).Create(&enum).Error; err != nil {
		return entities.Enum{}, err
	}
	return enum, nil
}

// Delete выполняет мягкое удаление Enum (устанавливает DeletedAt)
func (r *enumRepository) Delete(ctx context.Context, enum entities.Enum) error {
	if err := r.db.WithContext(ctx).Delete(&enum).Error; err != nil {
		return err
	}
	return nil
}

// Update обновляет существующую запись Enum
func (r *enumRepository) Update(ctx context.Context, enum entities.Enum) (entities.Enum, error) {
	if err := r.db.WithContext(ctx).Save(&enum).Error; err != nil {
		return entities.Enum{}, err
	}
	return enum, nil
}

// FindAll ищет все enum
func (r *enumRepository) FindAll(ctx context.Context) ([]entities.Enum, error) {
	var enums []entities.Enum
	if err := r.db.WithContext(ctx).Find(&enums).Error; err != nil {
		return nil, err
	}
	return enums, nil
}

// FindById ищет Enum по ID
func (r *enumRepository) FindById(ctx context.Context, id uint) (entities.Enum, error) {
	var enum entities.Enum
	if err := r.db.WithContext(ctx).First(&enum, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Enum{}, common.NewNotFoundError("enum not found by id")
		}
		return entities.Enum{}, err
	}
	return enum, nil
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
