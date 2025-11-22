package services

import (
	"context"
	"errors"
	"log/slog"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
)

const (
	cartServiceCode = "CART_SERVICE"
)

type CartService interface {
	BaseService[entities.Cart]

	Create(ctx context.Context, cart entities.Cart) (entities.Cart, error)
	Update(ctx context.Context, cart entities.Cart) (entities.Cart, error)
	Delete(ctx context.Context, cart entities.Cart) error

	GetByPersonID(ctx context.Context, personID uint) (entities.Cart, error)
}

type cartService struct {
	BaseServiceImpl[entities.Cart]
	cartRepo repositories.CartRepository
}

func NewCartService(
	cartRepo repositories.CartRepository,
) *cartService {
	return &cartService{
		BaseServiceImpl: *NewBaseServiceImpl(cartRepo),
		cartRepo:        cartRepo,
	}
}

// Create создает новую Cart с базовой валидацией
func (s *cartService) Create(ctx context.Context, cart entities.Cart) (entities.Cart, error) {
	// Создаем запись
	created, err := s.repo.Create(ctx, cart)
	if err != nil {
		slog.Error("Create cart failed", "error", err, "personID", cart.PersonID)
		return entities.Cart{}, appError.NewTechnicalError(err, cartServiceCode, err.Error())
	}
	slog.Info("Cart created", "personID", created.PersonID)
	return created, nil
}

// Update обновляет существующую Cart
func (s *cartService) Update(ctx context.Context, cart entities.Cart) (entities.Cart, error) {
	existing, err := s.repo.FindByID(ctx, cart.ID)
	if err != nil {
		return entities.Cart{}, err
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		slog.Error("Update cart failed", "error", err, "personID", cart.PersonID)
		return entities.Cart{}, err
	}
	return updated, nil
}

// Delete удаляет Cart (мягко или полностью)
func (s *cartService) Delete(ctx context.Context, cart entities.Cart) error {
	err := s.repo.Delete(ctx, cart)
	if err != nil {
		slog.Error("Failed to delete cart", "error", err, "id", cart.ID)
		return appError.NewTechnicalError(err, cartServiceCode, err.Error())
	}
	slog.Info("Deleted cart", "personID", cart.PersonID)
	return nil
}

// FindByCode ищет Cart по коду
func (s *cartService) GetByPersonID(ctx context.Context, personID uint) (entities.Cart, error) {
	cart, err := s.cartRepo.FindByPersonID(ctx, personID)
	if err != nil {
		slog.Error("Failed to find cart by code", "error", err, "personID", personID)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Cart{}, appError.NewLogicalError(err, cartServiceCode, err.Error())
		}
		return entities.Cart{}, appError.NewTechnicalError(err, cartServiceCode, err.Error())
	}
	return cart, nil
}
