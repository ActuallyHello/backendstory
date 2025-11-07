package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type BaseRepository[T entities.BaseEntity] interface {
	Create(ctx context.Context, entity T) (T, error)
	Update(ctx context.Context, entity T) (T, error)
	Delete(ctx context.Context, entity T) error

	FindAll(ctx context.Context) ([]T, error)
	FindByID(ctx context.Context, id uint) (T, error)
	FindWithSearchCriteria(ctx context.Context, criteria dto.SearchCriteria) ([]T, error)

	Count(ctx context.Context, criteria *dto.SearchCriteria) (int64, error)
}

type BaseRepositoryImpl[T entities.BaseEntity] struct {
	db *gorm.DB
}

func NewBaseRepositoryImpl[T entities.BaseEntity](db *gorm.DB) *BaseRepositoryImpl[T] {
	return &BaseRepositoryImpl[T]{
		db: db,
	}
}

// Create создает новую запись
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity T) (T, error) {
	if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Update обновляет существующую запись
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity T) (T, error) {
	if err := r.db.WithContext(ctx).Save(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Delete выполняет мягкое удаление
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, entity T) error {
	if err := r.db.WithContext(ctx).Delete(&entity).Error; err != nil {
		return err
	}
	return nil
}

// FindByID ищет запись по ID
func (r *BaseRepositoryImpl[T]) FindByID(ctx context.Context, id uint) (T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, common.NewNotFoundError("record not found by id")
		}
		return entity, err
	}
	return entity, nil
}

// FindAll ищет все записи
func (r *BaseRepositoryImpl[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// FindWithSearchCriteria ищет записи по критериям поиска
func (r *BaseRepositoryImpl[T]) FindWithSearchCriteria(ctx context.Context, criteria dto.SearchCriteria) ([]T, error) {
	var entities []T
	q := r.db.WithContext(ctx)
	queryCtx := common.BuildQuery(q, criteria)
	if err := queryCtx.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count возвращает количество записей по критериям
func (r *BaseRepositoryImpl[T]) Count(ctx context.Context, criteria *dto.SearchCriteria) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	if criteria != nil {
		query = common.BuildQuery(query, *criteria)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
