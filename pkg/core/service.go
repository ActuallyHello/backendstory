package core

import (
	"context"
	"errors"
)

const (
	serviceCodeSuffix = "_SERVICE"
)

// BaseService интерфейс для базовых операций сервиса
type BaseService[T BaseEntity] interface {
	GetRepo() BaseRepository[T]

	GetByID(ctx context.Context, id uint) (T, error)
	GetAll(ctx context.Context) ([]T, error)
	GetWithSearchCriteria(ctx context.Context, criteria SearchCriteria) ([]T, error)
}

// BaseServiceImpl базовая реализация сервиса
type BaseServiceImpl[T BaseEntity] struct {
	repo BaseRepository[T]
}

// NewBaseService создает новый базовый сервис
func NewBaseServiceImpl[T BaseEntity](
	repo BaseRepository[T],
) *BaseServiceImpl[T] {
	return &BaseServiceImpl[T]{
		repo: repo,
	}
}

func (s *BaseServiceImpl[T]) GetRepo() BaseRepository[T] {
	return s.repo
}

// GetByID получает сущность по ID
func (s *BaseServiceImpl[T]) GetByID(ctx context.Context, id uint) (T, error) {
	var empty T
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, &NotFoundError{}) {
			return empty, NewLogicalError(err, entity.TableName()+serviceCodeSuffix, err.Error())
		}
		return empty, NewTechnicalError(err, entity.TableName()+serviceCodeSuffix, err.Error())
	}
	return entity, nil
}

// GetAll получает все сущности
func (s *BaseServiceImpl[T]) GetAll(ctx context.Context) ([]T, error) {
	var entity T
	entities, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, NewTechnicalError(err, entity.TableName()+serviceCodeSuffix, "Ошибка при получении списка сущностей "+entity.TableName())
	}
	return entities, nil
}

// GetWithSearchCriteria ищет сущности по критериям
func (s *BaseServiceImpl[T]) GetWithSearchCriteria(ctx context.Context, criteria SearchCriteria) ([]T, error) {
	var entity T
	entities, err := s.repo.FindWithSearchCriteria(ctx, criteria)
	if err != nil {
		return nil, NewTechnicalError(err, entity.TableName()+serviceCodeSuffix, "Ошибка при получении по заданным параметрам сущностей "+entity.LocalTableName())
	}
	return entities, nil
}
