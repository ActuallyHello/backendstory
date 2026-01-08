package category

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	categoryServiceCode = "CATEGORY_SERVICE"
)

type CategoryService interface {
	core.BaseService[Category]

	Create(ctx context.Context, category Category) (Category, error)
	Update(ctx context.Context, category Category) (Category, error)
	Delete(ctx context.Context, category Category) error

	GetByCode(ctx context.Context, code string) (Category, error)
	GetByCategoryID(ctx context.Context, categoryID uint) ([]Category, error)
}

type categoryService struct {
	core.BaseServiceImpl[Category]
	categoryRepo CategoryRepository
}

func NewCategoryService(
	categoryRepo CategoryRepository,
) *categoryService {
	return &categoryService{
		BaseServiceImpl: *core.NewBaseServiceImpl(categoryRepo),
		categoryRepo:    categoryRepo,
	}
}

// Create создает новую Category с базовой валидацией
func (s *categoryService) Create(ctx context.Context, category Category) (Category, error) {
	// Проверка существования с таким кодом
	existing, err := s.GetByCode(ctx, category.Code)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return Category{}, err
	}
	if existing.ID > 0 {
		slog.Error("Category already exists!", "error", err, "code", category.Code)
		return Category{}, core.NewLogicalError(nil, categoryServiceCode, "Категория уже существует")
	}

	// Создаем запись
	created, err := s.GetRepo().Create(ctx, category)
	if err != nil {
		slog.Error("Create category failed", "error", err, "code", category.Code)
		return Category{}, core.NewTechnicalError(err, categoryServiceCode, "Ошибка при создании категории")
	}
	slog.Info("Category created", "code", created.Code)
	return created, nil
}

// Update обновляет существующую Category
func (s *categoryService) Update(ctx context.Context, category Category) (Category, error) {
	existing, err := s.GetRepo().FindByID(ctx, category.ID)
	if err != nil {
		return Category{}, err
	}

	updated, err := s.GetRepo().Update(ctx, existing)
	if err != nil {
		slog.Error("Update category failed", "error", err, "code", category.Code)
		return Category{}, core.NewTechnicalError(err, categoryServiceCode, "Ошибка при обновлении категории")
	}
	return updated, nil
}

// Delete удаляет Category (мягко или полностью)
func (s *categoryService) Delete(ctx context.Context, category Category) error {
	err := s.GetRepo().Delete(ctx, category)
	if err != nil {
		slog.Error("Failed to delete category", "error", err, "id", category.ID)
		return core.NewTechnicalError(err, categoryServiceCode, "Ошибка при удалении категории")
	}
	slog.Info("Deleted category", "code", category.Code)
	return nil
}

// FindByCode ищет Category по коду
func (s *categoryService) GetByCode(ctx context.Context, code string) (Category, error) {
	category, err := s.categoryRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find category by code", "error", err, "code", code)
		if errors.Is(err, &core.NotFoundError{}) {
			return Category{}, core.NewLogicalError(err, categoryServiceCode, err.Error())
		}
		return Category{}, core.NewTechnicalError(err, categoryServiceCode, "Ошибка при поиске категории по коду")
	}
	return category, nil
}

// FindByCategoryID ищет Category по category id
func (s *categoryService) GetByCategoryID(ctx context.Context, categoryID uint) ([]Category, error) {
	categories, err := s.categoryRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		slog.Error("Failed to find category by code", "error", err, "categoryID", categoryID)
		return nil, core.NewTechnicalError(err, categoryServiceCode, "Ошибка поиска категории по родителю")
	}
	return categories, nil
}
