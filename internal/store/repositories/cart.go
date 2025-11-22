package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type CartRepository interface {
	BaseRepository[entities.Cart]

	FindByPersonID(ctx context.Context, personID uint) (entities.Cart, error)
}

type cartRepository struct {
	BaseRepositoryImpl[entities.Cart]
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.Cart](db),
	}
}

// FindByPersonID ищет Cart по коду
func (r *cartRepository) FindByPersonID(ctx context.Context, personID uint) (entities.Cart, error) {
	var cart entities.Cart
	if err := r.db.WithContext(ctx).Where("PERSONID = ?", personID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Cart{}, common.NewNotFoundError("Cart not found by code")
		}
		return entities.Cart{}, err
	}
	return cart, nil
}
