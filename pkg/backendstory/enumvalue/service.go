package enumvalue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	enumValueServiceCode = "ENUMERATION_VALUE_SERVICE"
)

type EnumValueService interface {
	core.BaseService[EnumValue]

	Create(ctx context.Context, enumValue EnumValue) (EnumValue, error)
	Update(ctx context.Context, enumValue EnumValue) (EnumValue, error)
	Delete(ctx context.Context, enumValue EnumValue) error

	GetByEnumID(ctx context.Context, enumerationID uint) ([]EnumValue, error)
	GetByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (EnumValue, error)
}

type enumValueService struct {
	core.BaseServiceImpl[EnumValue]
	enumValueRepo EnumValueRepository
}

func NewEnumValueService(
	enumValueRepo EnumValueRepository,
) *enumValueService {
	return &enumValueService{
		BaseServiceImpl: *core.NewBaseServiceImpl(enumValueRepo),
		enumValueRepo:   enumValueRepo,
	}
}

// Create создает новую EnumValue с валидацией
func (s *enumValueService) Create(ctx context.Context, enumValue EnumValue) (EnumValue, error) {
	// Проверяем уникальность кода в рамках Enum
	existing, err := s.GetByCodeAndEnumID(ctx, enumValue.Code, enumValue.EnumID)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return EnumValue{}, err
	}
	if existing.ID > 0 {
		slog.Error("Enum value already exists for this enumeration!", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
		return EnumValue{}, core.NewLogicalError(
			nil,
			enumValueServiceCode,
			fmt.Sprintf("Значение перечеслимого типа с кодом %s уже существует!", enumValue.Code),
		)
	}

	// Создаем запись
	created, err := s.enumValueRepo.Create(ctx, enumValue)
	if err != nil {
		slog.Error("Create enumeration value failed", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, "Невозможно создать значение перечислимого типа")
	}
	slog.Info("Enum value created", "code", enumValue.Code, "enumID", enumValue.EnumID)
	return created, nil
}

// Update обновляет существующую EnumValue
func (s *enumValueService) Update(ctx context.Context, enumValue EnumValue) (EnumValue, error) {
	// Проверяем существование записи
	existing, err := s.GetByID(ctx, enumValue.ID)
	if err != nil {
		return EnumValue{}, err
	}

	// Если изменился код, проверяем уникальность
	if existing.Code != enumValue.Code {
		existingByCode, err := s.GetByCodeAndEnumID(ctx, enumValue.Code, enumValue.EnumID)
		if err != nil && errors.Is(err, &core.TechnicalError{}) {
			return EnumValue{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Enum value already exists for this enumeration!", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
			return EnumValue{}, core.NewLogicalError(
				err,
				enumValueServiceCode,
				fmt.Sprintf("Значение перечислимого с кодом %s уже существует!", enumValue.Code),
			)
		}
	}

	// Обновляем запись
	updated, err := s.enumValueRepo.Update(ctx, enumValue)
	if err != nil {
		slog.Error("Update enumeration value failed", "error", err, "id", enumValue.ID, "code", enumValue.Code)
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	slog.Info("Enum value updated", "id", enumValue.ID, "code", enumValue.Code)
	return updated, nil
}

// Delete удаляет EnumValue (мягко или полностью)
func (s *enumValueService) Delete(ctx context.Context, enumValue EnumValue) error {
	err := s.enumValueRepo.Delete(ctx, enumValue)
	if err != nil {
		slog.Error("Delete enumeration value failed", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
		return core.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	slog.Info("Deleted enumeration value", "code", enumValue.Code, "enumID", enumValue.EnumID)
	return nil
}

// GetByEnumID ищет все EnumValue по EnumID
func (s *enumValueService) GetByEnumID(ctx context.Context, enumerationID uint) ([]EnumValue, error) {
	values, err := s.enumValueRepo.FindByEnumID(ctx, enumerationID)
	if err != nil {
		slog.Error("Failed to find enumeration values by enumeration ID", "error", err, "enumID", enumerationID)
		return nil, core.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	return values, nil
}

// GetByCodeAndEnumID ищет EnumValue по коду и EnumID
func (s *enumValueService) GetByCodeAndEnumID(ctx context.Context, code string, enumerationID uint) (EnumValue, error) {
	// Используем репозиторий для поиска по коду, затем проверяем EnumID
	enumValue, err := s.enumValueRepo.FindByCodeAndEnumID(ctx, code, enumerationID)
	if err != nil {
		slog.Error("Failed to find enumeration value", "error", err, "code", code, "enumID", enumerationID)
		if errors.Is(err, &core.NotFoundError{}) {
			return EnumValue{}, core.NewLogicalError(err, enumValueServiceCode, err.Error())
		}
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, err.Error())
	}
	return enumValue, nil // Принадлежит другому Enum - считаем не найденной
}
