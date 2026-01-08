package orderitem

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	core.BaseRepository[OrderItem]

	FindByOrderID(ctx context.Context, statusID uint) ([]OrderItem, error)
}

type orderItemRepository struct {
	core.BaseRepositoryImpl[OrderItem]
}

func NewOrderItemRepository(db *gorm.DB) *orderItemRepository {
	return &orderItemRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[OrderItem](db),
	}
}

func (r *orderItemRepository) FindByOrderID(ctx context.Context, orderID uint) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := r.GetDB(ctx).Where("ORDERID = ?", orderID).Find(&orderItems).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("Не существует заказов с данным статусом")
		}
		return nil, err
	}
	return orderItems, nil
}
