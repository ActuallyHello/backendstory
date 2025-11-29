package purchase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/cart"
	cartitem "github.com/ActuallyHello/backendstory/pkg/backendstory/cart_item"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/person"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	purchaseServiceCode = "PURCHASE_SERVICE"
)

type PurchaseService interface {
	AddToCart(ctx context.Context, product product.Product, person person.Person, quantity uint) error
	Purchase(ctx context.Context, cartID uint) error
}

type purchaseService struct {
	txManager        core.TxManager
	cartService      cart.CartService
	cartItemService  cartitem.CartItemService
	productService   product.ProductService
	enumService      enum.EnumService
	enumValueService enumvalue.EnumValueService
}

func NewPurchaseService(
	cartSerivice cart.CartService,
	cartItemService cartitem.CartItemService,
	productService product.ProductService,
	enumService enum.EnumService,
	enumValueService enumvalue.EnumValueService,
) *purchaseService {
	return &purchaseService{
		cartService:      cartSerivice,
		cartItemService:  cartItemService,
		productService:   productService,
		enumService:      enumService,
		enumValueService: enumValueService,
	}
}

func (s *purchaseService) AddToCart(ctx context.Context, product product.Product, person person.Person, quantity uint) error {
	if product.Quantity < quantity {
		return core.NewLogicalError(
			nil,
			purchaseServiceCode,
			fmt.Sprintf("Невозможно добавить продукт %s в количестве %d шт. Максимум доступно %d шт.", product.Label, quantity, product.Quantity))
	}

	cart, err := s.cartService.GetByPersonID(ctx, person.ID)
	if err != nil {
		return core.NewTechnicalError(
			err,
			purchaseServiceCode,
			fmt.Sprintf("У пользователя с логином %s нет доступных корзин", person.UserLogin),
		)
	}

	cartItem := cartitem.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  quantity,
	}
	if _, err := s.cartItemService.Create(ctx, cartItem); err != nil {
		return core.NewTechnicalError(
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
		return core.NewTechnicalError(err, purchaseServiceCode, err.Error())
	}

	createdStatus, err := s.enumValueService.GetByCodeAndEnumID(ctx, "Created", status.ID)
	if err != nil {
		return core.NewLogicalError(err, purchaseServiceCode, err.Error())
	}
	pendingStatus, err := s.enumValueService.GetByCodeAndEnumID(ctx, "Pending", status.ID)
	if err != nil {
		return core.NewLogicalError(err, purchaseServiceCode, err.Error())
	}

	cartItems, err := s.cartItemService.GetByCartID(ctx, cartID)
	if err != nil {
		return core.NewTechnicalError(err, purchaseServiceCode, err.Error())
	}

	for _, cartItem := range cartItems {
		switch cartItem.StatusID {
		case createdStatus.ID:
			product, err := s.productService.GetByID(ctx, cartItem.ProductID)
			if err != nil {
				return err
			}

			err = s.txManager.Do(ctx, func(ctx context.Context) error {
				cartItem.StatusID = pendingStatus.ID
				if _, err := s.cartItemService.Update(ctx, cartItem); err != nil {
					return err
				}

				if product.Quantity < cartItem.Quantity {
					return core.NewLogicalError(nil, purchaseServiceCode, "Невозможно оформить товар! Текущее количество товара меньше чем запрашиваемые на покупку!")
				}
				product.Quantity = product.Quantity - cartItem.Quantity
				if _, err := s.productService.Update(ctx, product); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		case pendingStatus.ID:
			slog.Info("CartItem already in status wait!", "cartItem", cartItem.ID, "cart", cartID)
		}
	}
	return nil
}
