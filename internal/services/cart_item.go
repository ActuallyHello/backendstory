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
	cartItemServiceCode = "CART_ITEM_SERVICE"
)

type CartItemService interface {
	BaseService[entities.CartItem]

	Create(ctx context.Context, cartItem entities.CartItem) (entities.CartItem, error)
	Update(ctx context.Context, cartItem entities.CartItem) (entities.CartItem, error)
	Delete(ctx context.Context, cartItem entities.CartItem) error

	GetByCartID(ctx context.Context, cartID uint) ([]entities.CartItem, error)
}

type cartItemService struct {
	BaseServiceImpl[entities.CartItem]
	cartItemRepo repositories.CartItemRepository

	enumService      EnumService
	enumValueService EnumValueService
}

func NewCartItemService(
	cartItemRepo repositories.CartItemRepository,
	enumService EnumService,
	enumValueService EnumValueService,
) *cartItemService {
	return &cartItemService{
		BaseServiceImpl:  *NewBaseServiceImpl(cartItemRepo),
		cartItemRepo:     cartItemRepo,
		enumService:      enumService,
		enumValueService: enumValueService,
	}
}

// Create создает новую CartItem с базовой валидацией
func (s *cartItemService) Create(ctx context.Context, cartItem entities.CartItem) (entities.CartItem, error) {
	status, err := s.enumService.GetByCode(ctx, "CartItemStatus")
	if err != nil {
		return entities.CartItem{}, err
	}
	createdStatus, err := s.enumValueService.GetByCodeAndEnumID(ctx, "Created", status.ID)
	if err != nil {
		return entities.CartItem{}, err
	}
	cartItem.StatusID = createdStatus.ID

	// Создаем запись
	created, err := s.repo.Create(ctx, cartItem)
	if err != nil {
		slog.Error("Create cartItem failed", "error", err, "cartID", cartItem.CartID)
		return entities.CartItem{}, appError.NewTechnicalError(err, cartItemServiceCode, err.Error())
	}
	slog.Info("CartItem created", "cartID", created.CartID)
	return created, nil
}

// Update обновляет существующую CartItem
func (s *cartItemService) Update(ctx context.Context, cartItem entities.CartItem) (entities.CartItem, error) {
	if _, err := s.enumValueService.GetByID(ctx, cartItem.StatusID); err != nil {
		return entities.CartItem{}, err
	}

	updated, err := s.repo.Update(ctx, cartItem)
	if err != nil {
		slog.Error("Update cartItem failed", "error", err, "cartID", cartItem.CartID)
		return entities.CartItem{}, err
	}
	return updated, nil
}

// Delete удаляет CartItem (мягко или полностью)
func (s *cartItemService) Delete(ctx context.Context, cartItem entities.CartItem) error {
	err := s.repo.Delete(ctx, cartItem)
	if err != nil {
		slog.Error("Failed to delete cartItem", "error", err, "id", cartItem.ID)
		return appError.NewTechnicalError(err, cartItemServiceCode, err.Error())
	}
	slog.Info("Deleted cartItem", "cartID", cartItem.CartID)
	return nil
}

// FindByCode ищет CartItem по коду
func (s *cartItemService) GetByCartID(ctx context.Context, cartID uint) ([]entities.CartItem, error) {
	cartItems, err := s.cartItemRepo.FindByCartID(ctx, cartID)
	if err != nil {
		slog.Error("Failed to find cartItem by code", "error", err, "cartID", cartID)
		if errors.Is(err, &common.NotFoundError{}) {
			return nil, appError.NewLogicalError(err, cartItemServiceCode, err.Error())
		}
		return nil, appError.NewTechnicalError(err, cartItemServiceCode, err.Error())
	}
	return cartItems, nil
}
