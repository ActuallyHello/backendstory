package productmedia

import (
	"context"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	productMediaServiceCode = "PRODUCTMEDIA_SERVICE"
)

type ProductMediaService interface {
	core.BaseService[ProductMedia]

	Create(ctx context.Context, productMedia ProductMedia) (ProductMedia, error)
	Update(ctx context.Context, productMedia ProductMedia) (ProductMedia, error)
	Delete(ctx context.Context, productMedia ProductMedia) error

	GetByProductID(ctx context.Context, productID uint) ([]ProductMedia, error)
}

type productMediaService struct {
	core.BaseServiceImpl[ProductMedia]
	productMediaRepo ProductMediaRepository
}

func NewProductMediaService(
	productMediaRepo ProductMediaRepository,
) *productMediaService {
	return &productMediaService{
		BaseServiceImpl:  *core.NewBaseServiceImpl(productMediaRepo),
		productMediaRepo: productMediaRepo,
	}
}

// Create создает новую ProductMedia с базовой валидацией
func (s *productMediaService) Create(ctx context.Context, productMedia ProductMedia) (ProductMedia, error) {
	created, err := s.GetRepo().Create(ctx, productMedia)
	if err != nil {
		return ProductMedia{}, core.NewTechnicalError(err, productMediaServiceCode, err.Error())
	}
	return created, nil
}

// Update обновляет существующую ProductMedia
func (s *productMediaService) Update(ctx context.Context, productMedia ProductMedia) (ProductMedia, error) {
	existing, err := s.GetRepo().FindByID(ctx, productMedia.ID)
	if err != nil {
		return ProductMedia{}, err
	}

	updated, err := s.GetRepo().Update(ctx, existing)
	if err != nil {
		return ProductMedia{}, err
	}
	return updated, nil
}

// Delete удаляет ProductMedia (мягко или полностью)
func (s *productMediaService) Delete(ctx context.Context, productMedia ProductMedia) error {
	err := s.GetRepo().Delete(ctx, productMedia)
	if err != nil {
		return core.NewTechnicalError(err, productMediaServiceCode, err.Error())
	}
	return nil
}

// FindByProductID ищет ProductMedia по product id
func (s *productMediaService) GetByProductID(ctx context.Context, productID uint) ([]ProductMedia, error) {
	productMedia, err := s.productMediaRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, core.NewTechnicalError(err, productMediaServiceCode, err.Error())
	}
	return productMedia, nil
}
