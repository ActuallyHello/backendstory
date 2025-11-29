package person

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
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
		slog.Error("Person already exists for this user!", "error", err, "user_login", person.UserLogin)
		return Person{}, core.NewLogicalError(nil, personServiceCode, fmt.Sprintf("Person with user_login = %s already exists!", person.UserLogin))
	}

	// Создаем запись
	created, err := s.personRepo.Create(ctx, person)
	if err != nil {
		slog.Error("Create person failed", "error", err, "user_login", person.UserLogin, "phone", person.Phone)
		return Person{}, core.NewTechnicalError(err, personServiceCode, err.Error())
	}
	slog.Info("Person created", "firstname", created.Firstname, "lastname", created.Lastname, "user_login", created.UserLogin)
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
		slog.Error("Update person failed", "error", err, "personID", person.ID, "user_login", person.UserLogin)
		return Person{}, core.NewTechnicalError(err, personServiceCode, err.Error())
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
			slog.Error("Failed to soft delete person", "error", err, "personID", person.ID, "soft", soft)
			return err
		}
	} else {
		err := s.personRepo.Delete(ctx, person)
		if err != nil {
			slog.Error("Failed to delete person", "error", err, "personID", person.ID, "soft", soft)
			return core.NewTechnicalError(err, personServiceCode, err.Error())
		}
	}
	slog.Info("Deleted person", "personID", person.ID, "soft", soft)
	return nil
}

// FindByUserLogin ищет Person по UserLogin
func (s *personService) GetByUserLogin(ctx context.Context, userLogin string) (Person, error) {
	person, err := s.personRepo.FindByUserLogin(ctx, userLogin)
	if err != nil {
		slog.Error("Failed to find person by user ID", "error", err, "user_login", userLogin)
		if errors.Is(err, &common.NotFoundError{}) {
			return Person{}, core.NewLogicalError(err, personServiceCode, err.Error())
		}
		return Person{}, core.NewTechnicalError(err, personServiceCode, err.Error())
	}
	return person, nil
}
