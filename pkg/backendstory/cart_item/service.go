package cartitem

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
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
	GetByCartIDAndProductID(ctx context.Context, cartID, productID uint) (CartItem, error)
}

type cartItemService struct {
	core.BaseServiceImpl[CartItem]
	cartItemRepo     CartItemRepository
	enumService      enum.EnumService
	enumValueService enumvalue.EnumValueService
	productService   product.ProductService
}

func NewCartItemService(
	cartItemRepo CartItemRepository,
	enumService enum.EnumService,
	enumValueService enumvalue.EnumValueService,
	productService product.ProductService,
) *cartItemService {
	return &cartItemService{
		BaseServiceImpl:  *core.NewBaseServiceImpl(cartItemRepo),
		cartItemRepo:     cartItemRepo,
		enumService:      enumService,
		enumValueService: enumValueService,
		productService:   productService,
	}
}

func (s *cartItemService) Create(ctx context.Context, cartItem CartItem) (CartItem, error) {
	existing, err := s.GetByCartIDAndProductID(ctx, cartItem.CartID, cartItem.ProductID)
	if err != nil {
		if logicalErr, ok := err.(*core.LogicalError); !ok || !errors.Is(logicalErr.Err, &core.NotFoundError{}) {
			return CartItem{}, err
		}
	}
	if existing.ID > 0 {
		slog.Error("Create cartItem failed", "error", "Product already in cart", "cartID", cartItem.CartID)
		return CartItem{}, core.NewLogicalError(nil, cartItemServiceCode, "Данный товар уже существует в корзине!")
	}

	if err := s.checkProductQuantity(ctx, cartItem); err != nil {
		return CartItem{}, err
	}

	created, err := s.GetRepo().Create(ctx, cartItem)
	if err != nil {
		slog.Error("Create cartItem failed", "error", err, "cartID", cartItem.CartID)
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при создании элемента корзины")
	}
	slog.Info("CartItem created", "cartID", created.CartID)
	return created, nil
}

func (s *cartItemService) Update(ctx context.Context, cartItem CartItem) (CartItem, error) {
	if err := s.checkProductQuantity(ctx, cartItem); err != nil {
		return CartItem{}, err
	}

	updated, err := s.GetRepo().Update(ctx, cartItem)
	if err != nil {
		slog.Error("Update cartItem failed", "error", err, "cartID", cartItem.CartID)
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при обновлении элемента корзины")
	}
	return updated, nil
}

func (s *cartItemService) Delete(ctx context.Context, cartItem CartItem) error {
	err := s.GetRepo().Delete(ctx, cartItem)
	if err != nil {
		slog.Error("Failed to delete cartItem", "error", err, "id", cartItem.ID)
		return core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при удалении элемента корзины")
	}
	slog.Info("Deleted cartItem", "cartID", cartItem.CartID)
	return nil
}

func (s *cartItemService) GetByCartID(ctx context.Context, cartID uint) ([]CartItem, error) {
	cartItems, err := s.cartItemRepo.FindByCartID(ctx, cartID)
	if err != nil {
		slog.Error("Failed to find cartItem by cart", "error", err, "cartID", cartID)
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, cartItemServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при поиске элементов корзины")
	}
	return cartItems, nil
}

func (s *cartItemService) GetByCartIDAndProductID(ctx context.Context, cartID, productID uint) (CartItem, error) {
	cartItem, err := s.cartItemRepo.FindByCartIDAndProductID(ctx, cartID, productID)
	if err != nil {
		slog.Error("Failed to find cartItem by cart and product", "error", err, "cartID", cartID, "productID", productID)
		if errors.Is(err, &core.NotFoundError{}) {
			return CartItem{}, core.NewLogicalError(err, cartItemServiceCode, err.Error())
		}
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при поиске элементов корзины")
	}
	return cartItem, nil
}

func (s *cartItemService) checkProductQuantity(ctx context.Context, cartItem CartItem) error {
	product, err := s.productService.GetByID(ctx, cartItem.ProductID)
	if err != nil {
		return err
	}
	if product.Quantity < cartItem.Quantity {
		slog.Error("Create cartItem failed", "error", "Product quantity is less than cart item request", "cartID", cartItem.CartID)
		return core.NewLogicalError(nil, cartItemServiceCode, fmt.Sprintf("Невозможно создать элемент корзины! Текущее количество товара %s: %d", product.Label, product.Quantity))
	}
	return nil
}
