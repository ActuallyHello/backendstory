package core

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T BaseEntity] interface {
	GetDB(ctx context.Context) *gorm.DB

	Create(ctx context.Context, entity T) (T, error)
	Update(ctx context.Context, entity T) (T, error)
	Delete(ctx context.Context, entity T) error

	FindAll(ctx context.Context) ([]T, error)
	FindByID(ctx context.Context, id uint) (T, error)
	FindWithSearchCriteria(ctx context.Context, criteria SearchCriteria) ([]T, error)

	Count(ctx context.Context, criteria SearchCriteria) (int64, error)
}

type BaseRepositoryImpl[T BaseEntity] struct {
	db *gorm.DB
}

func NewBaseRepositoryImpl[T BaseEntity](db *gorm.DB) *BaseRepositoryImpl[T] {
	return &BaseRepositoryImpl[T]{
		db: db,
	}
}

func (r *BaseRepositoryImpl[T]) GetDB(ctx context.Context) *gorm.DB {
	// Пробуем получить транзакцию из контекста
	if tx, ok := ctx.Value(TxCtxKeyCode).(*gorm.DB); ok {
		return tx
	}
	// Если транзакции нет, используем обычное соединение
	return r.db.WithContext(ctx)
}

// Create создает новую запись
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity T) (T, error) {
	if err := r.GetDB(ctx).Create(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Update обновляет существующую запись
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity T) (T, error) {
	if err := r.GetDB(ctx).Save(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

// Delete выполняет мягкое удаление
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, entity T) error {
	if err := r.GetDB(ctx).Delete(&entity).Error; err != nil {
		return err
	}
	return nil
}

// FindByID ищет запись по ID
func (r *BaseRepositoryImpl[T]) FindByID(ctx context.Context, id uint) (T, error) {
	var entity T
	if err := r.GetDB(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, NewNotFoundError(entity.LocalTableName() + " с переданным ID не существует")
		}
		return entity, err
	}
	return entity, nil
}

// FindAll ищет все записи
func (r *BaseRepositoryImpl[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.GetDB(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// FindWithSearchCriteria ищет записи по критериям поиска
func (r *BaseRepositoryImpl[T]) FindWithSearchCriteria(ctx context.Context, criteria SearchCriteria) ([]T, error) {
	var entities []T
	q := r.GetDB(ctx)
	queryCtx := BuildQuery(q, criteria)

	if err := queryCtx.Debug().Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Count возвращает количество записей по критериям
func (r *BaseRepositoryImpl[T]) Count(ctx context.Context, criteria SearchCriteria) (int64, error) {
	var count int64
	query := r.GetDB(ctx).Model(new(T))

	query = BuildQuery(query, criteria)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
