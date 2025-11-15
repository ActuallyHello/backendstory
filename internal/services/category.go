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
	categoryServiceCode = "CATEGORY_SERVICE"
)

type CategoryService interface {
	BaseService[entities.Category]

	Create(ctx context.Context, category entities.Category) (entities.Category, error)
	Update(ctx context.Context, category entities.Category) (entities.Category, error)
	Delete(ctx context.Context, category entities.Category) error

	GetByCode(ctx context.Context, code string) (entities.Category, error)
	GetByCategoryID(ctx context.Context, categoryID uint) ([]entities.Category, error)
}

type categoryService struct {
	BaseServiceImpl[entities.Category]
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(
	categoryRepo repositories.CategoryRepository,
) *categoryService {
	return &categoryService{
		BaseServiceImpl: *NewBaseServiceImpl(categoryRepo),
		categoryRepo:    categoryRepo,
	}
}

// Create создает новую Category с базовой валидацией
func (s *categoryService) Create(ctx context.Context, category entities.Category) (entities.Category, error) {
	// Проверка существования с таким кодом
	existing, err := s.GetByCode(ctx, category.Code)
	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
		return entities.Category{}, err
	}
	if existing.ID > 0 {
		slog.Error("Category already exists!", "error", err, "code", category.Code)
		return entities.Category{}, appError.NewLogicalError(nil, categoryServiceCode, fmt.Sprintf("Category with code = %s already exists!", category.Code))
	}

	// Создаем запись
	created, err := s.repo.Create(ctx, category)
	if err != nil {
		slog.Error("Create category failed", "error", err, "code", category.Code)
		return entities.Category{}, appError.NewTechnicalError(err, categoryServiceCode, err.Error())
	}
	slog.Info("Category created", "code", created.Code)
	return created, nil
}

// Update обновляет существующую Category
func (s *categoryService) Update(ctx context.Context, category entities.Category) (entities.Category, error) {
	existing, err := s.repo.FindByID(ctx, category.ID)
	if err != nil {
		return entities.Category{}, err
	}

	if existing.Code != category.Code {
		existingByCode, err := s.GetByCode(ctx, category.Code)
		if err != nil && errors.Is(err, &appError.TechnicalError{}) {
			return entities.Category{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Category already exists!", "error", err, "code", category.Code)
			return entities.Category{}, appError.NewLogicalError(err, categoryServiceCode, fmt.Sprintf("Category with code = %s already exists!", category.Code))
		}
	}

	if category.Code != "" {
		existing.Code = category.Code
	}
	if category.Label != "" {
		existing.Label = category.Label
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		slog.Error("Update category failed", "error", err, "code", category.Code)
		return entities.Category{}, err
	}
	return updated, nil
}

// Delete удаляет Category (мягко или полностью)
func (s *categoryService) Delete(ctx context.Context, category entities.Category) error {
	err := s.repo.Delete(ctx, category)
	if err != nil {
		slog.Error("Failed to delete category", "error", err, "id", category.ID)
		return appError.NewTechnicalError(err, categoryServiceCode, err.Error())
	}
	slog.Info("Deleted category", "code", category.Code)
	return nil
}

// FindByCode ищет Category по коду
func (s *categoryService) GetByCode(ctx context.Context, code string) (entities.Category, error) {
	category, err := s.categoryRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find category by code", "error", err, "code", code)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Category{}, appError.NewLogicalError(err, categoryServiceCode, err.Error())
		}
		return entities.Category{}, appError.NewTechnicalError(err, categoryServiceCode, err.Error())
	}
	return category, nil
}

// FindByCategoryID ищет Category по category id
func (s *categoryService) GetByCategoryID(ctx context.Context, categoryID uint) ([]entities.Category, error) {
	categories, err := s.categoryRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		slog.Error("Failed to find category by code", "error", err, "categoryID", categoryID)
		return nil, appError.NewTechnicalError(err, categoryServiceCode, err.Error())
	}
	return categories, nil
}
