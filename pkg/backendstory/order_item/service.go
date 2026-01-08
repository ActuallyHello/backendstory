package orderitem

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	cartitem "github.com/ActuallyHello/backendstory/pkg/backendstory/cart_item"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	orderItemServiceCode = "ORDER_SERVICE"
)

type OrderItemService interface {
	core.BaseService[OrderItem]

	Create(ctx context.Context, orderItem OrderItem) (OrderItem, error)
	Update(ctx context.Context, orderItem OrderItem) (OrderItem, error)
	Delete(ctx context.Context, orderItem OrderItem) error

	Approve(ctx context.Context, orderItem OrderItem) (OrderItem, error)
	Cancel(ctx context.Context, orderItem OrderItem) (OrderItem, error)
	ChangeStatus(ctx context.Context, orderItem OrderItem, status string) (OrderItem, error)

	GetByOrderID(ctx context.Context, orderID uint) ([]OrderItem, error)
}

type orderItemService struct {
	core.BaseServiceImpl[OrderItem]
	orderItemRepo    OrderItemRepository
	txManager        core.TxManager
	enumService      enum.EnumService
	enumValueService enumvalue.EnumValueService
	productService   product.ProductService
	cartItemService  cartitem.CartItemService
}

func NewOrderItemService(
	orderItemRepo OrderItemRepository,
	txManager core.TxManager,
	enumService enum.EnumService,
	enumValueService enumvalue.EnumValueService,
	productService product.ProductService,
	cartItemService cartitem.CartItemService,
) *orderItemService {
	return &orderItemService{
		BaseServiceImpl:  *core.NewBaseServiceImpl(orderItemRepo),
		orderItemRepo:    orderItemRepo,
		txManager:        txManager,
		enumService:      enumService,
		enumValueService: enumValueService,
		productService:   productService,
		cartItemService:  cartItemService,
	}
}

func (s *orderItemService) Create(ctx context.Context, orderItem OrderItem) (OrderItem, error) {
	createdOrderItemStatusID, err := s.enumValueService.GetByCodeAndEnumCode(ctx, PendingOrderItemStatus, OrderItemStatus)
	if err != nil {
		return OrderItem{}, err
	}

	orderItem.StatusID = createdOrderItemStatusID.ID
	created, err := s.GetRepo().Create(ctx, orderItem)
	if err != nil {
		slog.Error("Create orderItem failed", "error", err, "orderID", orderItem.OrderID)
		return OrderItem{}, core.NewTechnicalError(err, orderItemServiceCode, "Ошибка при создании элемента заказа")
	}
	slog.Info("OrderItem created", "orderID", created.OrderID)
	return created, nil
}

func (s *orderItemService) Update(ctx context.Context, orderItem OrderItem) (OrderItem, error) {
	updated, err := s.GetRepo().Update(ctx, orderItem)
	if err != nil {
		slog.Error("Update orderItem failed", "error", err, "orderID", orderItem.OrderID)
		return OrderItem{}, core.NewTechnicalError(err, orderItemServiceCode, "Ошибка при обновлении элемента заказа")
	}
	slog.Info("OrderItem updated", "orderID", updated.OrderID)
	return updated, nil
}

func (s *orderItemService) ChangeStatus(ctx context.Context, orderItem OrderItem, status string) (OrderItem, error) {
	switch status {
	case ApprovedOrderItemStatus:
		return s.Approve(ctx, orderItem)
	case CancelledOrderItemStatus:
		return s.Cancel(ctx, orderItem)
	default:
		return OrderItem{}, core.NewLogicalError(nil, orderItemServiceCode, "Невозможно изменить статус заказа на "+status)
	}
}

func (s *orderItemService) Approve(ctx context.Context, orderItem OrderItem) (OrderItem, error) {
	var approvedOrderItem OrderItem
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		currentStatus, err := s.enumValueService.GetByID(ctx, orderItem.StatusID)
		if err != nil {
			return err
		}
		if slices.Contains([]string{CancelledOrderItemStatus}, currentStatus.Code) {
			return core.NewLogicalError(nil, orderItemServiceCode, fmt.Sprintf("Невозможно подтверджить заказ! Элемент заказа в статусе '%s'", currentStatus.Label))
		}

		approvedStatus, err := s.enumValueService.GetByCodeAndEnumCode(ctx, ApprovedOrderItemStatus, OrderItemStatus)
		if err != nil {
			return err
		}

		cartItem, err := s.cartItemService.GetByID(ctx, orderItem.CartItemID)
		if err != nil {
			return err
		}

		product, err := s.productService.GetByID(ctx, cartItem.ProductID)
		if err != nil {
			return err
		}

		if product.Quantity < cartItem.Quantity {
			return core.NewLogicalError(nil, orderItemServiceCode, fmt.Sprintf("Невозможно подтвердить элемент заказа! Текущее количество товара %s: %d", product.Label, product.Quantity))
		}

		product.Quantity = product.Quantity - cartItem.Quantity
		product, err = s.productService.Update(ctx, product)
		if err != nil {
			return err
		}

		orderItem.StatusID = approvedStatus.ID
		orderItem, err = s.Update(ctx, orderItem)
		if err != nil {
			return err
		}
		approvedOrderItem = orderItem

		return nil
	})
	return approvedOrderItem, err
}

func (s *orderItemService) Cancel(ctx context.Context, orderItem OrderItem) (OrderItem, error) {
	cancelStatus, err := s.enumValueService.GetByCodeAndEnumCode(ctx, CancelledOrderItemStatus, OrderItemStatus)
	if err != nil {
		return OrderItem{}, err
	}
	orderItem.StatusID = cancelStatus.ID
	return s.Update(ctx, orderItem)
}

func (s *orderItemService) Delete(ctx context.Context, orderItem OrderItem) error {
	err := s.GetRepo().Delete(ctx, orderItem)
	if err != nil {
		slog.Error("Failed to delete orderItem", "error", err, "id", orderItem.ID)
		return core.NewTechnicalError(err, orderItemServiceCode, "Ошибка при удалении элемента заказа")
	}
	slog.Info("Deleted orderItem", "orderID", orderItem.OrderID)
	return nil
}

func (s *orderItemService) GetByOrderID(ctx context.Context, orderID uint) ([]OrderItem, error) {
	orderItems, err := s.orderItemRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		slog.Error("Failed to find orderItem by order", "error", err, "orderID", orderID)
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, orderItemServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, orderItemServiceCode, "Ошибка при поиске элемента заказа")
	}
	return orderItems, nil
}
