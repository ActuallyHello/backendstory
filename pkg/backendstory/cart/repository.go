package cart

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type CartRepository interface {
	core.BaseRepository[Cart]

	FindByPersonID(ctx context.Context, personID uint) (Cart, error)
}

type cartRepository struct {
	core.BaseRepositoryImpl[Cart]
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[Cart](db),
	}
}

// FindByPersonID ищет Cart по коду
func (r *cartRepository) FindByPersonID(ctx context.Context, personID uint) (Cart, error) {
	var cart Cart
	if err := r.GetDB().WithContext(ctx).Where("PERSONID = ?", personID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Cart{}, core.NewNotFoundError("Cart not found by code")
		}
		return Cart{}, err
	}
	return cart, nil
}
