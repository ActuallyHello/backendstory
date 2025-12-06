package person

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type PersonRepository interface {
	core.BaseRepository[Person]

	FindByUserLogin(ctx context.Context, userLogin string) (Person, error)
}

type personRepository struct {
	core.BaseRepositoryImpl[Person]
}

func NewPersonRepository(db *gorm.DB) *personRepository {
	return &personRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[Person](db),
	}
}

// FindByUserID ищет Person по UserID (отношение 1:1)
func (r *personRepository) FindByUserLogin(ctx context.Context, userLogin string) (Person, error) {
	var person Person
	if err := r.GetDB().WithContext(ctx).Where("USERID = ?", userLogin).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Person{}, core.NewNotFoundError("person not found by user login")
		}
		return Person{}, err
	}
	return person, nil
}
