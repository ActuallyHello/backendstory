package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	productServiceCode = "PRODUCT_SERVICE"
)

type ProductService interface {
	core.BaseService[Product]

	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, product Product) error

	GetByCode(ctx context.Context, code string) (Product, error)
	GetByCategoryID(ctx context.Context, categoryID uint) ([]Product, error)
}

type productService struct {
	core.BaseServiceImpl[Product]
	productRepo ProductRepository

	enumService      enum.EnumService
	enumValueService enumvalue.EnumValueService
}

func NewProductService(
	productRepo ProductRepository,
	enumService enum.EnumService,
	enumValueService enumvalue.EnumValueService,
) *productService {
	return &productService{
		BaseServiceImpl:  *core.NewBaseServiceImpl(productRepo),
		productRepo:      productRepo,
		enumService:      enumService,
		enumValueService: enumValueService,
	}
}

// Create создает новую Product с базовой валидацией
func (s *productService) Create(ctx context.Context, product Product) (Product, error) {
	criteriaByCodeAndSku := core.SearchCriteria{
		Limit: 1,
		SearchConditions: []core.SearchCondition{
			{
				Field:     "CODE",
				Operation: core.OpEqual,
				Value:     product.Code,
			},
			{
				Field:     "SKU",
				Operation: core.OpEqual,
				Value:     product.Sku,
			},
		},
	}
	// Проверка существования с таким кодом
	products, err := s.GetWithSearchCriteria(ctx, criteriaByCodeAndSku)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return Product{}, err
	}
	if len(products) > 0 {
		slog.Error("Product already exists!", "error", err, "code", product.Code, "sku", product.Sku)
		return Product{}, core.NewLogicalError(nil, productServiceCode, fmt.Sprintf("Product with code = %s already exists!", product.Code))
	}

	// Создаем запись
	created, err := s.GetRepo().Create(ctx, product)
	if err != nil {
		slog.Error("Create product failed", "error", err, "code", product.Code, "sku", product.Sku)
		return Product{}, core.NewTechnicalError(err, productServiceCode, err.Error())
	}
	slog.Info("Product created", "code", created.Code, "sku", product.Sku)
	return created, nil
}

// Update обновляет существующую Product
func (s *productService) Update(ctx context.Context, product Product) (Product, error) {
	existing, err := s.GetRepo().FindByID(ctx, product.ID)
	if err != nil {
		return Product{}, err
	}

	if existing.Code != product.Code {
		existingByCode, err := s.GetByCode(ctx, product.Code)
		if err != nil && errors.Is(err, &core.TechnicalError{}) {
			return Product{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Product already exists!", "error", err, "code", product.Code)
			return Product{}, core.NewLogicalError(err, productServiceCode, fmt.Sprintf("Product with code = %s already exists!", product.Code))
		}
	}

	if product.Code != "" {
		existing.Code = product.Code
	}
	if product.Label != "" {
		existing.Label = product.Label
	}

	updated, err := s.GetRepo().Update(ctx, existing)
	if err != nil {
		slog.Error("Update product failed", "error", err, "code", product.Code)
		return Product{}, err
	}
	return updated, nil
}

// Delete удаляет Product (мягко или полностью)
func (s *productService) Delete(ctx context.Context, product Product) error {
	err := s.GetRepo().Delete(ctx, product)
	if err != nil {
		slog.Error("Failed to delete product", "error", err, "id", product.ID)
		return core.NewTechnicalError(err, productServiceCode, err.Error())
	}
	slog.Info("Deleted product", "code", product.Code)
	return nil
}

// FindByCode ищет Product по коду
func (s *productService) GetByCode(ctx context.Context, code string) (Product, error) {
	product, err := s.productRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find product by code", "error", err, "code", code)
		if errors.Is(err, &core.NotFoundError{}) {
			return Product{}, core.NewLogicalError(err, productServiceCode, err.Error())
		}
		return Product{}, core.NewTechnicalError(err, productServiceCode, err.Error())
	}
	return product, nil
}

// FindByProductID ищет Product по category id
func (s *productService) GetByCategoryID(ctx context.Context, categoryID uint) ([]Product, error) {
	products, err := s.productRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		slog.Error("Failed to find product by code", "error", err, "categoryID", categoryID)
		return nil, core.NewTechnicalError(err, productServiceCode, err.Error())
	}
	return products, nil
}
