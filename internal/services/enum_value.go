package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
)

const (
	enumValueServiceCode = "ENUMERATION_VALUE_SERVICE"
)

type EnumValueService interface {
	Create(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error)
	Update(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error)
	Delete(ctx context.Context, enumValue entities.EnumValue) error

	GetAll(ctx context.Context) ([]entities.EnumValue, error)
	GetById(ctx context.Context, id uint) (entities.EnumValue, error)
	GetByEnumID(ctx context.Context, enumerationID uint) ([]entities.EnumValue, error)
	GetByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (entities.EnumValue, error)
}

type enumValueService struct {
	enumValueRepo repositories.EnumValueRepository
}

func NewEnumValueService(
	enumValueRepo repositories.EnumValueRepository,
) *enumValueService {
	return &enumValueService{
		enumValueRepo: enumValueRepo,
	}
}

// Create создает новую EnumValue с валидацией
func (s *enumValueService) Create(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error) {
	// Проверяем уникальность кода в рамках Enum
	existing, err := s.GetByCodeAndEnumID(ctx, enumValue.Code, enumValue.EnumID)
	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
		return entities.EnumValue{}, err
	}
	if existing.ID > 0 {
		slog.Error("Enum value already exists for this enumeration!", "error", err, "code", enumValue.Code, "enumeration_id", enumValue.EnumID)
		return entities.EnumValue{}, appError.NewLogicalError(
			nil,
			enumValueServiceCode,
			fmt.Sprintf("Enum value with code = %s already exists for enumeration ID = %d!", enumValue.Code, enumValue.EnumID),
		)
	}

	// Создаем запись
	created, err := s.enumValueRepo.Create(ctx, enumValue)
	if err != nil {
		slog.Error("Create enumeration value failed", "error", err, "code", enumValue.Code, "enumeration_id", enumValue.EnumID)
		return entities.EnumValue{}, appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	slog.Info("Enum value created", "code", enumValue.Code, "enumeration_id", enumValue.EnumID)
	return created, nil
}

// Update обновляет существующую EnumValue
func (s *enumValueService) Update(ctx context.Context, enumValue entities.EnumValue) (entities.EnumValue, error) {
	// Проверяем существование записи
	existing, err := s.GetById(ctx, enumValue.ID)
	if err != nil {
		return entities.EnumValue{}, err
	}

	// Если изменился код, проверяем уникальность
	if existing.Code != enumValue.Code {
		existingByCode, err := s.GetByCodeAndEnumID(ctx, enumValue.Code, enumValue.EnumID)
		if err != nil && errors.Is(err, &appError.TechnicalError{}) {
			return entities.EnumValue{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Enum value already exists for this enumeration!", "error", err, "code", enumValue.Code, "enumeration_id", enumValue.EnumID)
			return entities.EnumValue{}, appError.NewLogicalError(
				err,
				enumValueServiceCode,
				fmt.Sprintf("Enum value with code = %s already exists for enumeration ID = %d!", enumValue.Code, enumValue.EnumID),
			)
		}
	}

	// Обновляем запись
	updated, err := s.enumValueRepo.Update(ctx, enumValue)
	if err != nil {
		slog.Error("Update enumeration value failed", "error", err, "id", enumValue.ID, "code", enumValue.Code)
		return entities.EnumValue{}, appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	slog.Info("Enum value updated", "id", enumValue.ID, "code", enumValue.Code)
	return updated, nil
}

// Delete удаляет EnumValue (мягко или полностью)
func (s *enumValueService) Delete(ctx context.Context, enumValue entities.EnumValue) error {
	err := s.enumValueRepo.Delete(ctx, enumValue)
	if err != nil {
		slog.Error("Delete enumeration value failed", "error", err, "code", enumValue.Code, "enumerationID", enumValue.EnumID)
		return appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	slog.Info("Deleted enumeration value", "code", enumValue.Code, "enumerationID", enumValue.EnumID)
	return nil
}

// GetById ищет EnumValue по ID
func (s *enumValueService) GetAll(ctx context.Context) ([]entities.EnumValue, error) {
	values, err := s.enumValueRepo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to find enumeration values", "error", err)
		return nil, appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	return values, nil
}

// GetById ищет EnumValue по ID
func (s *enumValueService) GetById(ctx context.Context, id uint) (entities.EnumValue, error) {
	value, err := s.enumValueRepo.FindById(ctx, id)
	if err != nil {
		slog.Error("Failed to find enumeration value by ID", "error", err, "id", id)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.EnumValue{}, appError.NewLogicalError(err, enumValueServiceCode, fmt.Sprintf("enumeration value with id=%d not found!", id))
		}
		return entities.EnumValue{}, appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	return value, nil
}

// GetByEnumID ищет все EnumValue по EnumID
func (s *enumValueService) GetByEnumID(ctx context.Context, enumerationID uint) ([]entities.EnumValue, error) {
	values, err := s.enumValueRepo.FindByEnumID(ctx, enumerationID)
	if err != nil {
		slog.Error("Failed to find enumeration values by enumeration ID", "error", err, "enumeration_id", enumerationID)
		return nil, appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	return values, nil
}

// GetByCodeAndEnumID ищет EnumValue по коду и EnumID
func (s *enumValueService) GetByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (entities.EnumValue, error) {
	// Используем репозиторий для поиска по коду, затем проверяем EnumID
	enumValue, err := s.enumValueRepo.FindByCodeAndEnumID(ctx, code, enumerationID)
	if err != nil {
		slog.Error("Failed to find enumeration value", "error", err, "code", code, "enumerationID", enumerationID)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.EnumValue{}, appError.NewLogicalError(err, enumValueServiceCode, err.Error())
		}
		return entities.EnumValue{}, appError.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	return enumValue, nil // Принадлежит другому Enum - считаем не найденной
}
