package enumvalue

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
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

	GetByEnumID(ctx context.Context, enumID uint) ([]EnumValue, error)
	GetByCodeAndEnumID(ctx context.Context, code string, enumID uint) (EnumValue, error)
	GetByCodeAndEnumCode(ctx context.Context, code, enumCode string) (EnumValue, error)
}

type enumValueService struct {
	core.BaseServiceImpl[EnumValue]
	enumValueRepo EnumValueRepository
	enumService   enum.EnumService
}

func NewEnumValueService(
	enumValueRepo EnumValueRepository,
	enumService enum.EnumService,
) *enumValueService {
	return &enumValueService{
		BaseServiceImpl: *core.NewBaseServiceImpl(enumValueRepo),
		enumValueRepo:   enumValueRepo,
		enumService:     enumService,
	}
}

func (s *enumValueService) Create(ctx context.Context, enumValue EnumValue) (EnumValue, error) {
	existing, err := s.GetByCodeAndEnumID(ctx, enumValue.Code, enumValue.EnumID)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return EnumValue{}, err
	}
	if existing.ID > 0 {
		slog.Error("Enum value already exists for this enumeration!", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
		return EnumValue{}, core.NewLogicalError(
			nil,
			enumValueServiceCode,
			"Значение перечислимого типа уже существует",
		)
	}

	created, err := s.enumValueRepo.Create(ctx, enumValue)
	if err != nil {
		slog.Error("Create enumeration value failed", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, "Невозможно создать значение перечислимого типа")
	}
	slog.Info("Enum value created", "code", enumValue.Code, "enumID", enumValue.EnumID)
	return created, nil
}

func (s *enumValueService) Update(ctx context.Context, enumValue EnumValue) (EnumValue, error) {
	_, err := s.GetByID(ctx, enumValue.ID)
	if err != nil {
		return EnumValue{}, err
	}

	updated, err := s.enumValueRepo.Update(ctx, enumValue)
	if err != nil {
		slog.Error("Update enumeration value failed", "error", err, "id", enumValue.ID, "code", enumValue.Code)
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, "Ошибка при обновлении значения перечислимого типа")
	}
	slog.Info("Enum value updated", "id", enumValue.ID, "code", enumValue.Code)
	return updated, nil
}

func (s *enumValueService) Delete(ctx context.Context, enumValue EnumValue) error {
	err := s.enumValueRepo.Delete(ctx, enumValue)
	if err != nil {
		slog.Error("Delete enumeration value failed", "error", err, "code", enumValue.Code, "enumID", enumValue.EnumID)
		return core.NewTechnicalError(err, enumValueServiceCode, "Ошибка при удалении значения перечислимого типа")
	}
	slog.Info("Deleted enumeration value", "code", enumValue.Code, "enumID", enumValue.EnumID)
	return nil
}

func (s *enumValueService) GetByEnumID(ctx context.Context, enumID uint) ([]EnumValue, error) {
	values, err := s.enumValueRepo.FindByEnumID(ctx, enumID)
	if err != nil {
		slog.Error("Failed to find enumeration values by enumeration ID", "error", err, "enumID", enumID)
		return nil, core.NewTechnicalError(err, enumValueServiceCode, "Ошибка при получении значения перечисления по родителю")
	}
	return values, nil
}

func (s *enumValueService) GetByCodeAndEnumID(ctx context.Context, code string, enumID uint) (EnumValue, error) {
	enumValue, err := s.enumValueRepo.FindByCodeAndEnumID(ctx, code, enumID)
	if err != nil {
		slog.Error("Failed to find enumeration value", "error", err, "code", code, "enumID", enumID)
		if errors.Is(err, &core.NotFoundError{}) {
			return EnumValue{}, core.NewLogicalError(err, enumValueServiceCode, err.Error())
		}
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, "Ошибка при поиске значения перечисления по коду и родителю")
	}
	return enumValue, nil
}

func (s *enumValueService) GetByCodeAndEnumCode(ctx context.Context, code, enumCode string) (EnumValue, error) {
	enum, err := s.enumService.GetByCode(ctx, enumCode)
	if err != nil {
		return EnumValue{}, err
	}

	enumValue, err := s.enumValueRepo.FindByCodeAndEnumID(ctx, code, enum.ID)
	if err != nil {
		slog.Error("Failed to find enumeration value", "error", err, "code", code, "enum code", enumCode)
		if errors.Is(err, &core.NotFoundError{}) {
			return EnumValue{}, core.NewLogicalError(err, enumValueServiceCode, err.Error())
		}
		return EnumValue{}, core.NewTechnicalError(err, enumValueServiceCode, "Ошибка при поиске значения перечисления по коду и родителю")
	}
	return enumValue, nil
}
