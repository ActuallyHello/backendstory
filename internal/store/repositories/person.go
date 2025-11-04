package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type PersonRepository interface {
	Create(ctx context.Context, person entities.Person) (entities.Person, error)
	Update(ctx context.Context, person entities.Person) (entities.Person, error)
	Delete(ctx context.Context, person entities.Person, soft bool) error

	FindAll(ctx context.Context) ([]entities.Person, error)
	FindById(ctx context.Context, id uint) (entities.Person, error)
	FindByUserLogin(ctx context.Context, userLogin string) (entities.Person, error)
}

type personRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *personRepository {
	return &personRepository{db: db}
}

// Create создает новую запись Person
func (r *personRepository) Create(ctx context.Context, person entities.Person) (entities.Person, error) {
	if err := r.db.WithContext(ctx).Create(&person).Error; err != nil {
		return entities.Person{}, err
	}
	return person, nil
}

// Delete выполняет удаление Person
func (r *personRepository) Delete(ctx context.Context, person entities.Person, soft bool) error {
	if soft {
		person.DeletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		_, err := r.Update(ctx, person)
		if err != nil {
			return err
		}
		return nil
	} else {
		if err := r.db.WithContext(ctx).Delete(&person).Error; err != nil {
			return err
		}
		return nil
	}
}

// Update обновляет данные Person
func (r *personRepository) Update(ctx context.Context, person entities.Person) (entities.Person, error) {
	if err := r.db.WithContext(ctx).Save(&person).Error; err != nil {
		return entities.Person{}, err
	}
	return person, nil
}

// FindAll ищет всех Person
func (r *personRepository) FindAll(ctx context.Context) ([]entities.Person, error) {
	var persons []entities.Person
	if err := r.db.WithContext(ctx).Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

// FindById ищет Person по ID
func (r *personRepository) FindById(ctx context.Context, id uint) (entities.Person, error) {
	var person entities.Person
	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Person{}, common.NewNotFoundError("person not found by id")
		}
		return entities.Person{}, err
	}
	return person, nil
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

// FindByPhone ищет Person по частичному совпадению номера телефона
// func (r *personRepository) FindByPhoneLike(ctx context.Context, phone string) ([]entities.Person, error) {
// 	var persons []entities.Person
// 	phonePattern := "%" + phone + "%"

// 	if err := r.db.WithContext(ctx).
// 		Where("PHONE LIKE ?", phonePattern).
// 		Find(&persons).Error; err != nil {
// 		slog.Error("FindByPhone failed", "error", err, "phone", phone)
// 		return nil, err
// 	}
// 	return persons, nil
// }

// FindByNames ищет Person по частичному совпадению имени и/или фамилии
// func (r *personRepository) FindByNames(ctx context.Context, firstname, lastname string) ([]entities.Person, error) {
// 	if firstname == "" && lastname == "" {
// 		return []entities.Person{}, nil
// 	}

// 	var persons []entities.Person

// 	query := r.db.WithContext(ctx).Model(&entities.Person{})

// 	conditions := make([]clause.Expression, 0)

// 	if firstname != "" {
// 		conditions = append(conditions, clause.Like{
// 			Column: "FIRSTNAME",
// 			Value:  "%" + firstname + "%",
// 		})
// 	}

// 	if lastname != "" {
// 		conditions = append(conditions, clause.Like{
// 			Column: "LASTNAME",
// 			Value:  "%" + lastname + "%",
// 		})
// 	}

// 	// Объединяем условия через AND
// 	query = query.Clauses(clause.And(conditions...))

// 	if err := query.Find(&persons).Error; err != nil {
// 		slog.Error("FindByNames failed", "error", err, "firstname", firstname, "lastname", lastname)
// 		return nil, err
// 	}

// 	return persons, nil
// }
