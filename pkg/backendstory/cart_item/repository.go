package cartitem

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	core.BaseRepository[CartItem]

	FindByCartID(ctx context.Context, cartID uint) ([]CartItem, error)
}

type cartItemRepository struct {
	core.BaseRepositoryImpl[CartItem]
}

func NewCartItemRepository(db *gorm.DB) *cartItemRepository {
	return &cartItemRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[CartItem](db),
	}
}

// FindByCartID ищет CartItem по Cart
func (r *cartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]CartItem, error) {
	var cartItems []CartItem
	if err := r.GetDB().WithContext(ctx).Where("CARTID = ?", cartID).Find(&cartItems).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("CartItem not found by code")
		}
		return nil, err
	}
	return cartItems, nil
}
