package services

import (
	"context"
	"log/slog"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
)

const (
	productMediaServiceCode = "PRODUCTMEDIA_SERVICE"
)

type ProductMediaService interface {
	BaseService[entities.ProductMedia]

	Create(ctx context.Context, productMedia entities.ProductMedia) (entities.ProductMedia, error)
	Update(ctx context.Context, productMedia entities.ProductMedia) (entities.ProductMedia, error)
	Delete(ctx context.Context, productMedia entities.ProductMedia) error

	GetByProductID(ctx context.Context, productID uint) ([]entities.ProductMedia, error)
}

type productMediaService struct {
	BaseServiceImpl[entities.ProductMedia]
	productMediaRepo repositories.ProductMediaRepository
}

func NewProductMediaService(
	productMediaRepo repositories.ProductMediaRepository,
) *productMediaService {
	return &productMediaService{
		BaseServiceImpl:  *NewBaseServiceImpl(productMediaRepo),
		productMediaRepo: productMediaRepo,
	}
}

// Create создает новую ProductMedia с базовой валидацией
func (s *productMediaService) Create(ctx context.Context, productMedia entities.ProductMedia) (entities.ProductMedia, error) {
	// Создаем запись
	created, err := s.repo.Create(ctx, productMedia)
	if err != nil {
		slog.Error("Create productMedia failed", "error", err, "productId", productMedia.ProductID)
		return entities.ProductMedia{}, appError.NewTechnicalError(err, productMediaServiceCode, err.Error())
	}
	slog.Info("ProductMedia created", "productId", created.ProductID)
	return created, nil
}

// Update обновляет существующую ProductMedia
func (s *productMediaService) Update(ctx context.Context, productMedia entities.ProductMedia) (entities.ProductMedia, error) {
	existing, err := s.repo.FindByID(ctx, productMedia.ID)
	if err != nil {
		return entities.ProductMedia{}, err
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		slog.Error("Update productMedia failed", "error", err, "productId", productMedia.ProductID)
		return entities.ProductMedia{}, err
	}
	return updated, nil
}

// Delete удаляет ProductMedia (мягко или полностью)
func (s *productMediaService) Delete(ctx context.Context, productMedia entities.ProductMedia) error {
	err := s.repo.Delete(ctx, productMedia)
	if err != nil {
		slog.Error("Failed to delete productMedia", "error", err, "id", productMedia.ID)
		return appError.NewTechnicalError(err, productMediaServiceCode, err.Error())
	}
	slog.Info("Deleted productMedia", "productId", productMedia.ProductID)
	return nil
}

// FindByProductID ищет ProductMedia по product id
func (s *productMediaService) GetByProductID(ctx context.Context, productID uint) ([]entities.ProductMedia, error) {
	productMedia, err := s.productMediaRepo.FindByProductID(ctx, productID)
	if err != nil {
		slog.Error("Failed to find productMedia by code", "error", err, "productId", productID)
		return nil, appError.NewTechnicalError(err, productMediaServiceCode, err.Error())
	}
	return productMedia, nil
}
