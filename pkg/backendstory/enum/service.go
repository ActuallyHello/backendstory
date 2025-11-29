package enum

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	enumServiceCode = "ENUMERATION_SERVICE"
)

type EnumService interface {
	core.BaseService[Enum]

	Create(ctx context.Context, enum Enum) (Enum, error)
	Update(ctx context.Context, enum Enum) (Enum, error)
	Delete(ctx context.Context, enum Enum) error

	GetByCode(ctx context.Context, code string) (Enum, error)
}

type enumService struct {
	core.BaseServiceImpl[Enum]
	enumRepo EnumRepository
}

func NewEnumService(
	enumRepo EnumRepository,
) *enumService {
	return &enumService{
		BaseServiceImpl: *core.NewBaseServiceImpl(enumRepo),
		enumRepo:        enumRepo,
	}
}

func (s *enumService) Create(ctx context.Context, enum Enum) (Enum, error) {
	existing, err := s.GetByCode(ctx, enum.Code)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return Enum{}, err
	}
	if existing.ID > 0 {
		slog.Error("Enum already exists!", "error", err, "code", enum.Code)
		return Enum{}, core.NewLogicalError(nil, enumServiceCode, "перечисление уже существует")
	}

	created, err := s.GetRepo().Create(ctx, enum)
	if err != nil {
		slog.Error("Create enum failed", "error", err, "code", enum.Code)
		return Enum{}, core.NewTechnicalError(err, enumServiceCode, "ошибка при создании перечисления")
	}
	slog.Info("Enum created", "code", created.Code)
	return created, nil
}

func (s *enumService) Update(ctx context.Context, enum Enum) (Enum, error) {
	existing, err := s.GetRepo().FindByID(ctx, enum.ID)
	if err != nil {
		return Enum{}, err
	}

	if existing.Code != enum.Code {
		existingByCode, err := s.GetByCode(ctx, enum.Code)
		if err != nil && errors.Is(err, &core.TechnicalError{}) {
			return Enum{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Enum already exists!", "error", err, "code", enum.Code)
			return Enum{}, core.NewLogicalError(err, enumServiceCode, "перечесление уже существует")
		}
	}

	existing.Label = enum.Label
	existing.Code = enum.Code

	updated, err := s.GetRepo().Update(ctx, existing)
	if err != nil {
		slog.Error("Update enum failed", "error", err, "code", enum.Code)
		return Enum{}, err
	}
	return updated, nil
}

func (s *enumService) Delete(ctx context.Context, enum Enum) error {
	err := s.GetRepo().Delete(ctx, enum)
	if err != nil {
		slog.Error("Failed to delete enum", "error", err, "id", enum.ID)
		return core.NewTechnicalError(err, enumServiceCode, "ошибка в удалении перечисления")
	}
	slog.Info("Deleted enum", "code", enum.Code)
	return nil
}

func (s *enumService) GetByCode(ctx context.Context, code string) (Enum, error) {
	enum, err := s.enumRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find enum by code", "error", err, "code", code)
		if errors.Is(err, &common.NotFoundError{}) {
			return Enum{}, core.NewLogicalError(err, enumServiceCode, err.Error())
		}
		return Enum{}, core.NewTechnicalError(err, enumServiceCode, "ошибка при получении перечисления по коду")
	}
	return enum, nil
}
