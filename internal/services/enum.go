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
	enumServiceCode = "ENUMERATION_SERVICE"
)

type EnumService interface {
	Create(ctx context.Context, enum entities.Enum) (entities.Enum, error)
	Update(ctx context.Context, enum entities.Enum) (entities.Enum, error)
	Delete(ctx context.Context, enum entities.Enum) error

	GetAll(ctx context.Context) ([]entities.Enum, error)
	GetById(ctx context.Context, id uint) (entities.Enum, error)
	GetByCode(ctx context.Context, code string) (entities.Enum, error)
}

type enumService struct {
	enumRepo repositories.EnumRepository
}

func NewEnumService(
	enumRepo repositories.EnumRepository,
) *enumService {
	return &enumService{
		enumRepo: enumRepo,
	}
}

// Create создает новую Enum с базовой валидацией
func (s *enumService) Create(ctx context.Context, enum entities.Enum) (entities.Enum, error) {
	// Проверка существования с таким кодом
	existing, err := s.GetByCode(ctx, enum.Code)
	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
		return entities.Enum{}, err
	}
	if existing.ID > 0 {
		slog.Error("Enum already exists!", "error", err, "code", enum.Code)
		return entities.Enum{}, appError.NewLogicalError(nil, enumServiceCode, fmt.Sprintf("Enum with code = %s already exists!", enum.Code))
	}

	// Создаем запись
	created, err := s.enumRepo.Create(ctx, enum)
	if err != nil {
		slog.Error("Create enum failed", "error", err, "code", enum.Code)
		return entities.Enum{}, appError.NewTechnicalError(err, enumServiceCode, err.Error())
	}
	slog.Info("Enum created", "code", created.Code)
	return created, nil
}

// Update обновляет существующую Enum
func (s *enumService) Update(ctx context.Context, enum entities.Enum) (entities.Enum, error) {
	existing, err := s.enumRepo.FindById(ctx, enum.ID)
	if err != nil {
		return entities.Enum{}, err
	}

	if existing.Code != enum.Code {
		existingByCode, err := s.GetByCode(ctx, enum.Code)
		if err != nil && errors.Is(err, &appError.TechnicalError{}) {
			return entities.Enum{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Enum already exists!", "error", err, "code", enum.Code)
			return entities.Enum{}, appError.NewLogicalError(err, enumServiceCode, fmt.Sprintf("Enum with code = %s already exists!", enum.Code))
		}
	}

	if enum.Code != "" {
		existing.Code = enum.Code
	}
	if enum.Label != "" {
		existing.Label = enum.Label
	}

	updated, err := s.enumRepo.Update(ctx, existing)
	if err != nil {
		slog.Error("Update enum failed", "error", err, "code", enum.Code)
		return entities.Enum{}, err
	}
	return updated, nil
}

// Delete удаляет Enum (мягко или полностью)
func (s *enumService) Delete(ctx context.Context, enum entities.Enum) error {
	err := s.enumRepo.Delete(ctx, enum)
	if err != nil {
		slog.Error("Failed to delete enum", "error", err, "id", enum.ID)
		return appError.NewTechnicalError(err, enumServiceCode, err.Error())
	}
	slog.Info("Deleted enum", "code", enum.Code)
	return nil
}

// FindAll ищет все Enum
func (s *enumService) GetAll(ctx context.Context) ([]entities.Enum, error) {
	enums, err := s.enumRepo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to find enums", "error", err)
		return nil, appError.NewTechnicalError(err, enumServiceCode, err.Error())
	}
	return enums, nil
}

// FindById ищет Enum по ID
func (s *enumService) GetById(ctx context.Context, id uint) (entities.Enum, error) {
	enum, err := s.enumRepo.FindById(ctx, id)
	if err != nil {
		slog.Error("Failed to find enum by ID", "error", err, "id", id)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Enum{}, appError.NewLogicalError(err, enumServiceCode, err.Error())
		}
		return entities.Enum{}, appError.NewTechnicalError(err, enumServiceCode, err.Error())
	}
	return enum, nil
}

// FindByCode ищет Enum по коду
func (s *enumService) GetByCode(ctx context.Context, code string) (entities.Enum, error) {
	enum, err := s.enumRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find enum by code", "error", err, "code", code)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Enum{}, appError.NewLogicalError(err, enumServiceCode, err.Error())
		}
		return entities.Enum{}, appError.NewTechnicalError(err, enumServiceCode, err.Error())
	}
	return enum, nil
}
