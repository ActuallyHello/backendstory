package category

import (
	"context"
	"errors"

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
		return Category{}, core.NewLogicalError(nil, categoryServiceCode, "Категория уже существует")
	}

	// Создаем запись
	created, err := s.GetRepo().Create(ctx, category)
	if err != nil {
		return Category{}, core.NewTechnicalError(err, categoryServiceCode, "Ошибка при создании категории")
	}
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
		return Category{}, core.NewTechnicalError(err, categoryServiceCode, "Ошибка при обновлении категории")
	}
	return updated, nil
}

// Delete удаляет Category (мягко или полностью)
func (s *categoryService) Delete(ctx context.Context, category Category) error {
	err := s.GetRepo().Delete(ctx, category)
	if err != nil {
		return core.NewTechnicalError(err, categoryServiceCode, "Ошибка при удалении категории")
	}
	return nil
}

// FindByCode ищет Category по коду
func (s *categoryService) GetByCode(ctx context.Context, code string) (Category, error) {
	category, err := s.categoryRepo.FindByCode(ctx, code)
	if err != nil {
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
		return nil, core.NewTechnicalError(err, categoryServiceCode, "Ошибка поиска категории по родителю")
	}
	return categories, nil
}
