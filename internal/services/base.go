package services

import (
	"context"
	"errors"
	"log/slog"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
)

const (
	serviceCodeSuffix = "_SERVICE"
)

// BaseService интерфейс для базовых операций сервиса
type BaseService[T entities.BaseEntity] interface {
	GetByID(ctx context.Context, id uint) (T, error)
	GetAll(ctx context.Context) ([]T, error)
	GetWithSearchCriteria(ctx context.Context, criteria dto.SearchCriteria) ([]T, error)
}

// BaseServiceImpl базовая реализация сервиса
type BaseServiceImpl[T entities.BaseEntity] struct {
	repo repositories.BaseRepository[T]
}

// NewBaseService создает новый базовый сервис
func NewBaseServiceImpl[T entities.BaseEntity](
	repo repositories.BaseRepository[T],
) *BaseServiceImpl[T] {
	return &BaseServiceImpl[T]{
		repo: repo,
	}
}

// GetByID получает сущность по ID
func (s *BaseServiceImpl[T]) GetByID(ctx context.Context, id uint) (T, error) {
	var zero T
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		slog.Error("GetByID failed", "error", err, "id", id, "entity", entity.TableName())
		if errors.Is(err, &common.NotFoundError{}) {
			return zero, appError.NewLogicalError(err, entity.TableName()+serviceCodeSuffix, "entity not found")
		}
		return zero, appError.NewTechnicalError(err, entity.TableName()+serviceCodeSuffix, err.Error())
	}
	return entity, nil
}

// GetAll получает все сущности
func (s *BaseServiceImpl[T]) GetAll(ctx context.Context) ([]T, error) {
	entities, err := s.repo.FindAll(ctx)
	if err != nil {
		slog.Error("GetAll failed", "error", err)
		return nil, appError.NewTechnicalError(err, T(*new(T)).TableName()+serviceCodeSuffix, err.Error())
	}
	return entities, nil
}

// GetWithSearchCriteria ищет сущности по критериям
func (s *BaseServiceImpl[T]) GetWithSearchCriteria(ctx context.Context, criteria dto.SearchCriteria) ([]T, error) {
	entities, err := s.repo.FindWithSearchCriteria(ctx, criteria)
	if err != nil {
		slog.Error("GetWithSearchCriteria failed", "error", err)
		return nil, appError.NewTechnicalError(err, T(*new(T)).TableName()+serviceCodeSuffix, err.Error())
	}
	return entities, nil
}
