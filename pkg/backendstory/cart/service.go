package cart

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	cartServiceCode = "CART_SERVICE"
)

type CartService interface {
	core.BaseService[Cart]

	Create(ctx context.Context, cart Cart) (Cart, error)
	Delete(ctx context.Context, cart Cart) error

	GetByPersonID(ctx context.Context, personID uint) (Cart, error)
}

type cartService struct {
	core.BaseServiceImpl[Cart]
	cartRepo CartRepository
}

func NewCartService(
	cartRepo CartRepository,
) *cartService {
	return &cartService{
		BaseServiceImpl: *core.NewBaseServiceImpl(cartRepo),
		cartRepo:        cartRepo,
	}
}

// Create создает новую Cart с базовой валидацией
func (s *cartService) Create(ctx context.Context, cart Cart) (Cart, error) {
	// Создаем запись
	created, err := s.GetRepo().Create(ctx, cart)
	if err != nil {
		slog.Error("Create cart failed", "error", err, "personID", cart.PersonID)
		return Cart{}, core.NewTechnicalError(err, cartServiceCode, "Ошибка при создании корзины")
	}
	slog.Info("Cart created", "personID", created.PersonID)
	return created, nil
}

// Delete удаляет Cart (мягко или полностью)
func (s *cartService) Delete(ctx context.Context, cart Cart) error {
	err := s.GetRepo().Delete(ctx, cart)
	if err != nil {
		slog.Error("Failed to delete cart", "error", err, "id", cart.ID)
		return core.NewTechnicalError(err, cartServiceCode, "Ошибка при удалении пользователя")
	}
	slog.Info("Deleted cart", "personID", cart.PersonID)
	return nil
}

// FindByCode ищет Cart по коду
func (s *cartService) GetByPersonID(ctx context.Context, personID uint) (Cart, error) {
	cart, err := s.cartRepo.FindByPersonID(ctx, personID)
	if err != nil {
		slog.Error("Failed to find cart by code", "error", err, "personID", personID)
		if errors.Is(err, &core.NotFoundError{}) {
			return Cart{}, core.NewLogicalError(err, cartServiceCode, err.Error())
		}
		return Cart{}, core.NewTechnicalError(err, cartServiceCode, "Ошибка при поиске корзины у клиента")
	}
	return cart, nil
}
