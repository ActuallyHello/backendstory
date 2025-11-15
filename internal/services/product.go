package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
)

const (
	productServiceCode = "PRODUCT_SERVICE"
)

type ProductService interface {
	BaseService[entities.Product]

	Create(ctx context.Context, product entities.Product) (entities.Product, error)
	Update(ctx context.Context, product entities.Product) (entities.Product, error)
	Delete(ctx context.Context, product entities.Product) error

	GetByCode(ctx context.Context, code string) (entities.Product, error)
	GetByCategoryID(ctx context.Context, categoryID uint) ([]entities.Product, error)
}

type productService struct {
	BaseServiceImpl[entities.Product]
	productRepo repositories.ProductRepository
}

func NewProductService(
	productRepo repositories.ProductRepository,
) *productService {
	return &productService{
		BaseServiceImpl: *NewBaseServiceImpl(productRepo),
		productRepo:     productRepo,
	}
}

// Create создает новую Product с базовой валидацией
func (s *productService) Create(ctx context.Context, product entities.Product) (entities.Product, error) {
	criteriaByCodeAndSku := dto.SearchCriteria{
		Limit: 1,
		SearchConditions: []dto.SearchCondition{
			dto.SearchCondition{
				Field:     "CODE",
				Operation: dto.OpEqual,
				Value:     product.Code,
			},
			dto.SearchCondition{
				Field:     "SKU",
				Operation: dto.OpEqual,
				Value:     product.Sku,
			},
		},
	}
	// Проверка существования с таким кодом
	products, err := s.GetWithSearchCriteria(ctx, criteriaByCodeAndSku)
	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
		return entities.Product{}, err
	}
	if len(products) > 0 {
		slog.Error("Product already exists!", "error", err, "code", product.Code, "sku", product.Sku)
		return entities.Product{}, appError.NewLogicalError(nil, productServiceCode, fmt.Sprintf("Product with code = %s already exists!", product.Code))
	}

	// Создаем запись
	created, err := s.repo.Create(ctx, product)
	if err != nil {
		slog.Error("Create product failed", "error", err, "code", product.Code, "sku", product.Sku)
		return entities.Product{}, appError.NewTechnicalError(err, productServiceCode, err.Error())
	}
	slog.Info("Product created", "code", created.Code, "sku", product.Sku)
	return created, nil
}

// Update обновляет существующую Product
func (s *productService) Update(ctx context.Context, product entities.Product) (entities.Product, error) {
	existing, err := s.repo.FindByID(ctx, product.ID)
	if err != nil {
		return entities.Product{}, err
	}

	if existing.Code != product.Code {
		existingByCode, err := s.GetByCode(ctx, product.Code)
		if err != nil && errors.Is(err, &appError.TechnicalError{}) {
			return entities.Product{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Product already exists!", "error", err, "code", product.Code)
			return entities.Product{}, appError.NewLogicalError(err, productServiceCode, fmt.Sprintf("Product with code = %s already exists!", product.Code))
		}
	}

	if product.Code != "" {
		existing.Code = product.Code
	}
	if product.Label != "" {
		existing.Label = product.Label
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		slog.Error("Update product failed", "error", err, "code", product.Code)
		return entities.Product{}, err
	}
	return updated, nil
}

// Delete удаляет Product (мягко или полностью)
func (s *productService) Delete(ctx context.Context, product entities.Product) error {
	err := s.repo.Delete(ctx, product)
	if err != nil {
		slog.Error("Failed to delete product", "error", err, "id", product.ID)
		return appError.NewTechnicalError(err, productServiceCode, err.Error())
	}
	slog.Info("Deleted product", "code", product.Code)
	return nil
}

// FindByCode ищет Product по коду
func (s *productService) GetByCode(ctx context.Context, code string) (entities.Product, error) {
	product, err := s.productRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find product by code", "error", err, "code", code)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Product{}, appError.NewLogicalError(err, productServiceCode, err.Error())
		}
		return entities.Product{}, appError.NewTechnicalError(err, productServiceCode, err.Error())
	}
	return product, nil
}

// FindByProductID ищет Product по category id
func (s *productService) GetByCategoryID(ctx context.Context, categoryID uint) ([]entities.Product, error) {
	products, err := s.productRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		slog.Error("Failed to find product by code", "error", err, "categoryID", categoryID)
		return nil, appError.NewTechnicalError(err, productServiceCode, err.Error())
	}
	return products, nil
}
