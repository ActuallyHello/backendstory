package person

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	personServiceCode = "PERSON_SERVICE"
)

type PersonService interface {
	core.BaseService[Person]

	Create(ctx context.Context, person Person) (Person, error)
	Update(ctx context.Context, person Person) (Person, error)
	Delete(ctx context.Context, person Person, soft bool) error

	GetByUserLogin(ctx context.Context, userLogin string) (Person, error)
}

type personService struct {
	core.BaseServiceImpl[Person]
	personRepo PersonRepository
}

func NewPersonService(
	personRepo PersonRepository,
) *personService {
	return &personService{
		BaseServiceImpl: *core.NewBaseServiceImpl(personRepo),
		personRepo:      personRepo,
	}
}

// Create создает новую запись Person
func (s *personService) Create(ctx context.Context, person Person) (Person, error) {
	// Проверка существования Person с таким UserLogin
	existingByUserLogin, err := s.personRepo.FindByUserLogin(ctx, person.UserLogin)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return Person{}, err
	}
	if existingByUserLogin.ID > 0 {
		return Person{}, core.NewLogicalError(nil, personServiceCode, "Клиент уже существует")
	}

	// Создаем запись
	created, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return Person{}, core.NewTechnicalError(err, personServiceCode, "Ошибка при создании клиента")
	}
	return created, nil
}

// Update обновляет существующую запись Person
func (s *personService) Update(ctx context.Context, person Person) (Person, error) {
	// Проверяем существование Person
	_, err := s.personRepo.FindByID(ctx, person.ID)
	if err != nil {
		return Person{}, err
	}

	updated, err := s.personRepo.Update(ctx, person)
	if err != nil {
		return Person{}, core.NewTechnicalError(err, personServiceCode, "Ошибка при обновлении клиента")
	}
	return updated, nil
}

// Delete удаляет Person (мягко или полностью)
func (s *personService) Delete(ctx context.Context, person Person, soft bool) error {
	if soft {
		person.DeletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if _, err := s.Update(ctx, person); err != nil {
			return err
		}
	} else {
		err := s.personRepo.Delete(ctx, person)
		if err != nil {
			return core.NewTechnicalError(err, personServiceCode, "Ошибка при удалении клиента")
		}
	}
	return nil
}

// FindByUserLogin ищет Person по UserLogin
func (s *personService) GetByUserLogin(ctx context.Context, userLogin string) (Person, error) {
	person, err := s.personRepo.FindByUserLogin(ctx, userLogin)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return Person{}, core.NewLogicalError(err, personServiceCode, err.Error())
		}
		return Person{}, core.NewTechnicalError(err, personServiceCode, "Ошибка при поиске клеинта по логину пользователя")
	}
	return person, nil
}
