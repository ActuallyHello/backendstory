package cartitem

import (
	"context"
	"errors"
	"fmt"

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
		return CartItem{}, core.NewLogicalError(nil, cartItemServiceCode, "Данный товар уже существует в корзине!")
	}

	if err := s.checkProduct(ctx, cartItem); err != nil {
		return CartItem{}, err
	}

	created, err := s.GetRepo().Create(ctx, cartItem)
	if err != nil {
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при создании элемента корзины")
	}
	return created, nil
}

func (s *cartItemService) Update(ctx context.Context, cartItem CartItem) (CartItem, error) {
	if err := s.checkProduct(ctx, cartItem); err != nil {
		return CartItem{}, err
	}

	updated, err := s.GetRepo().Update(ctx, cartItem)
	if err != nil {
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при обновлении элемента корзины")
	}
	return updated, nil
}

func (s *cartItemService) Delete(ctx context.Context, cartItem CartItem) error {
	err := s.GetRepo().Delete(ctx, cartItem)
	if err != nil {
		return core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при удалении элемента корзины")
	}
	return nil
}

func (s *cartItemService) GetByCartID(ctx context.Context, cartID uint) ([]CartItem, error) {
	cartItems, err := s.cartItemRepo.FindByCartID(ctx, cartID)
	if err != nil {
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
		if errors.Is(err, &core.NotFoundError{}) {
			return CartItem{}, core.NewLogicalError(err, cartItemServiceCode, err.Error())
		}
		return CartItem{}, core.NewTechnicalError(err, cartItemServiceCode, "Ошибка при поиске элементов корзины")
	}
	return cartItem, nil
}

func (s *cartItemService) checkProduct(ctx context.Context, cartItem CartItem) error {
	product, err := s.productService.GetByID(ctx, cartItem.ProductID)
	if err != nil {
		return err
	}
	if err := s.checkProductQuantity(cartItem, product); err != nil {
		return err
	}
	if err := s.checkProductStatus(ctx, product); err != nil {
		return err
	}
	return nil
}

func (s *cartItemService) checkProductQuantity(cartItem CartItem, product product.Product) error {
	if product.Quantity < cartItem.Quantity {
		return core.NewLogicalError(nil, cartItemServiceCode, fmt.Sprintf("Невозможно создать элемент корзины! Текущее количество товара %s: %d", product.Label, product.Quantity))
	}
	return nil
}

func (s *cartItemService) checkProductStatus(ctx context.Context, checkProduct product.Product) error {
	currentProductStatus, err := s.enumValueService.GetByID(ctx, checkProduct.StatusID)
	if err != nil {
		return err
	}
	if currentProductStatus.Code != product.AvailableProductStatus {
		return core.NewLogicalError(nil, cartItemServiceCode, "Товар недоступен для добавления в корзину")
	}
	return nil
}
