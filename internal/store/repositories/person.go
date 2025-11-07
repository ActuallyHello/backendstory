package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type PersonRepository interface {
	BaseRepository[entities.Person]

	FindByUserLogin(ctx context.Context, userLogin string) (entities.Person, error)
}

type personRepository struct {
	BaseRepositoryImpl[entities.Person]
}

func NewPersonRepository(db *gorm.DB) *personRepository {
	return &personRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.Person](db),
	}
}

// FindByUserID ищет Person по UserID (отношение 1:1)
func (r *personRepository) FindByUserLogin(ctx context.Context, userLogin string) (entities.Person, error) {
	var person entities.Person
	if err := r.db.WithContext(ctx).Where("USERID = ?", userLogin).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Person{}, common.NewNotFoundError("person not found by user login")
		}
		return entities.Person{}, err
	}
	return person, nil
}
