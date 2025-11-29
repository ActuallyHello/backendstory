package cartitem

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	cartItemServiceCode = "CART_ITEM_SERVICE"
)

type CartItemService interface {
	core.BaseService[CartItem]

	Create(ctx context.Context, cartItem CartItem) (CartItem, error)
	Update(ctx context.Context, cartItem CartItem) (CartItem, error)
	Delete(ctx context.Context, cartItem CartItem) error

	GetByCartID(ctx context.Context, cartID uint) ([]CartItem, error)
}

type cartItemService struct {
	core.BaseServiceImpl[CartItem]
	cartItemRepo CartItemRepository

	enumService      enum.EnumService
	enumValueService enumvalue.EnumValueService
}

func NewCartItemService(
	cartItemRepo CartItemRepository,
	enumService enum.EnumService,
	enumValueService enumvalue.EnumValueService,
) *cartItemService {
	return &cartItemService{
		BaseServiceImpl:  *core.NewBaseServiceImpl(cartItemRepo),
		cartItemRepo:     cartItemRepo,
		enumService:      enumService,
		enumValueService: enumValueService,
	}
}

// Create создает новую CartItem с базовой валидацией
func (s *cartItemService) Create(ctx context.Context, cartItem CartItem) (CartItem, error) {
	status, err := s.enumService.GetByCode(ctx, "CartItemStatus")
	if err != nil {
		return CartItem{}, err
	}
	createdStatus, err := s.enumValueService.GetByCodeAndEnumID(ctx, "Created", status.ID)
	if err != nil {
		return CartItem{}, err
	}
	cartItem.StatusID = createdStatus.ID

	// Создаем запись
	created, err := s.GetRepo().Create(ctx, cartItem)
	if err != nil {
		slog.Error("Create cartItem failed", "error", err, "cartID", cartItem.CartID)
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, err.Error())
	}
	slog.Info("CartItem created", "cartID", created.CartID)
	return created, nil
}

// Update обновляет существующую CartItem
func (s *cartItemService) Update(ctx context.Context, cartItem CartItem) (CartItem, error) {
	if _, err := s.enumValueService.GetByID(ctx, cartItem.StatusID); err != nil {
		return CartItem{}, err
	}

	updated, err := s.GetRepo().Update(ctx, cartItem)
	if err != nil {
		slog.Error("Update cartItem failed", "error", err, "cartID", cartItem.CartID)
		return CartItem{}, err
	}
	return updated, nil
}

// Delete удаляет CartItem (мягко или полностью)
func (s *cartItemService) Delete(ctx context.Context, cartItem CartItem) error {
	err := s.GetRepo().Delete(ctx, cartItem)
	if err != nil {
		slog.Error("Failed to delete cartItem", "error", err, "id", cartItem.ID)
		return core.NewTechnicalError(err, cartItemServiceCode, err.Error())
	}
	slog.Info("Deleted cartItem", "cartID", cartItem.CartID)
	return nil
}

// FindByCode ищет CartItem по коду
func (s *cartItemService) GetByCartID(ctx context.Context, cartID uint) ([]CartItem, error) {
	cartItems, err := s.cartItemRepo.FindByCartID(ctx, cartID)
	if err != nil {
		slog.Error("Failed to find cartItem by code", "error", err, "cartID", cartID)
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, cartItemServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, cartItemServiceCode, err.Error())
	}
	return cartItems, nil
}
