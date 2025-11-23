package services

import (
	"context"
	"fmt"
	"log/slog"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

const (
	purchaseServiceCode = "PURCHASE_SERVICE"
)

type PurchaseService interface {
	AddToCart(ctx context.Context, product entities.Product, person entities.Person, quantity uint) error
	Purchase(ctx context.Context, cartID uint) error
}

type purchaseService struct {
	cartService      CartService
	cartItemService  CartItemService
	productService   ProductService
	enumService      EnumService
	enumValueService EnumValueService
}

func NewPurchaseService(
	cartSerivice CartService,
	cartItemService CartItemService,
	productService ProductService,
	enumService EnumService,
	enumValueService EnumValueService,
) *purchaseService {
	return &purchaseService{
		cartService:      cartSerivice,
		cartItemService:  cartItemService,
		productService:   productService,
		enumService:      enumService,
		enumValueService: enumValueService,
	}
}

func (s *purchaseService) AddToCart(ctx context.Context, product entities.Product, person entities.Person, quantity uint) error {
	if product.Quantity < quantity {
		return appErr.NewLogicalError(
			nil,
			purchaseServiceCode,
			fmt.Sprintf("Невозможно добавить продукт %s в количестве %d шт. Максимум доступно %d шт.", product.Label, quantity, product.Quantity))
	}

	cart, err := s.cartService.GetByPersonID(ctx, person.ID)
	if err != nil {
		return appErr.NewTechnicalError(
			err,
			purchaseServiceCode,
			fmt.Sprintf("У пользователя с логином %s нет доступных корзин", person.UserLogin),
		)
	}

	cartItem := entities.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  quantity,
	}
	if _, err := s.cartItemService.Create(ctx, cartItem); err != nil {
		return appErr.NewTechnicalError(
			err,
			purchaseServiceCode,
			fmt.Sprintf("Ошибка во время добавления продукта %s в корзину", product.Label),
		)
	}

	return nil
}

func (s *purchaseService) Purchase(ctx context.Context, cartID uint) error {
	status, err := s.enumService.GetByCode(ctx, "CartItemStatus")
	if err != nil {
		return appErr.NewTechnicalError(err, purchaseServiceCode, err.Error())
	}

	createdStatus, err := s.enumValueService.GetByCodeAndEnumID(ctx, "Created", status.ID)
	if err != nil {
		return appErr.NewLogicalError(err, purchaseServiceCode, err.Error())
	}
	pendingStatus, err := s.enumValueService.GetByCodeAndEnumID(ctx, "Pending", status.ID)
	if err != nil {
		return appErr.NewLogicalError(err, purchaseServiceCode, err.Error())
	}

	cartItems, err := s.cartItemService.GetByCartID(ctx, cartID)
	if err != nil {
		return appErr.NewTechnicalError(err, purchaseServiceCode, err.Error())
	}

	for _, cartItem := range cartItems {
		switch cartItem.StatusID {
		case createdStatus.ID:
			cartItem.StatusID = pendingStatus.ID
			if _, err := s.cartItemService.Update(ctx, cartItem); err != nil {
				return err
			}
			product, err := s.productService.GetByID(ctx, cartItem.ProductID)
			if err != nil {
				return err
			}
			if product.Quantity < cartItem.Quantity {
				return appErr.NewLogicalError(nil, purchaseServiceCode, "Невозможно оформить товар! Текущее количество товара меньше чем запрашиваемые на покупку!")
			}
			product.Quantity = product.Quantity - cartItem.Quantity
			if _, err := s.productService.Update(ctx, product); err != nil {
				return err
			}
		case pendingStatus.ID:
			slog.Info("CartItem already in status wait!", "cartItem", cartItem.ID, "cart", cartID)
		}
	}
	return nil
}
