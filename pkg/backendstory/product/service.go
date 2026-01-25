package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	Delete(ctx context.Context, product Product, soft bool) error

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

func (s *productService) Create(ctx context.Context, product Product) (Product, error) {
	exists, err := s.isProductExists(ctx, product)
	if err != nil {
		return Product{}, err
	}
	if exists {
		return Product{}, core.NewLogicalError(nil, productServiceCode, "Продукт уже существует")
	}

	status, err := s.enumValueService.GetByCodeAndEnumCode(ctx, AvailableProductStatus, ProductStatus)
	if err != nil {
		return Product{}, err
	}

	product.StatusID = status.ID
	created, err := s.GetRepo().Create(ctx, product)
	if err != nil {
		return Product{}, core.NewTechnicalError(err, productServiceCode, err.Error())
	}
	return created, nil
}

func (s *productService) Update(ctx context.Context, product Product) (Product, error) {
	updated, err := s.GetRepo().Update(ctx, product)
	if err != nil {
		return Product{}, core.NewTechnicalError(err, productServiceCode, "Ошибка при обновлении продукта")
	}
	return updated, nil
}

func (s *productService) Delete(ctx context.Context, product Product, soft bool) error {
	var err error
	if soft {
		product.DeletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		_, err = s.GetRepo().Update(ctx, product)
	} else {
		err = s.GetRepo().Delete(ctx, product)
	}

	if err != nil {
		return core.NewTechnicalError(err, productServiceCode, "Ошибка при удалении продукта")
	}
	return nil
}

func (s *productService) GetByCode(ctx context.Context, code string) (Product, error) {
	product, err := s.productRepo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return Product{}, core.NewLogicalError(err, productServiceCode, err.Error())
		}
		return Product{}, core.NewTechnicalError(err, productServiceCode, "Ошибка при поиске продукта по переданному коду")
	}
	return product, nil
}

func (s *productService) GetByCategoryID(ctx context.Context, categoryID uint) ([]Product, error) {
	products, err := s.productRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, core.NewTechnicalError(err, productServiceCode, err.Error())
	}
	return products, nil
}

func (s *productService) isProductExists(ctx context.Context, product Product) (bool, error) {
	conditions := []core.SearchCondition{}
	conditions = append(conditions, core.SearchCondition{
		Field:     "CODE",
		Operation: core.OpEqual,
		Value:     product.Code,
	})
	conditions = append(conditions, core.SearchCondition{
		Field:     "SKU",
		Operation: core.OpEqual,
		Value:     product.Sku,
	})
	criteria := core.SearchCriteria{
		Limit:            1,
		SearchConditions: conditions,
	}

	products, err := s.GetWithSearchCriteria(ctx, criteria)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return false, err
	}
	if len(products) > 0 {
		return true, nil
	}
	return false, nil
}
