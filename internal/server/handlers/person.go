package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/server/handlers/common"
	"github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/go-playground/validator/v10"
)

const (
	personHandlerCode = "PERSON_HANDLER"
)

type PersonHandler struct {
	validate      *validator.Validate
	personService services.PersonService
}

func NewPersonHandler(
	personService services.PersonService,
) *PersonHandler {
	return &PersonHandler{
		validate:      validator.New(),
		personService: personService,
	}
}

// CreatePerson создает новую персону
// @Summary Создать персону
// @Description Создает новую запись персоны в системе
// @Tags persons
// @Accept json
// @Produce json
// @Param request body dto.CreatePersonRequest true "Данные для создания персоны"
// @Success 201 {object} dto.PersonDTO "Созданная персона"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 409 {object} dto.ErrorResponse "Конфликт данных"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons [post]
func (h *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, personHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, personHandlerCode, err.Error()), details)
		return
	}

	person := entities.Person{
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Phone:     req.Phone,
		UserLogin: req.UserLogin,
	}
	person, err := h.personService.Create(ctx, person)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.toPersonDTO(person))
}

// GetAll возвращает персон
// @Summary Получить персон
// @Description Возвращает персону по ее идентификатору
// @Tags persons
// @Produce json
// @Success 200 {object} dto.PersonDTO "Найденная персона"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} dto.ErrorResponse "Персона не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons [get]
func (h *PersonHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	persons, err := h.personService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.PersonDTO, 0, len(persons))
	for _, person := range persons {
		dtos = append(dtos, h.toPersonDTO(person))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает персону по ID
// @Summary Получить персону по ID
// @Description Возвращает персону по ее идентификатору
// @Tags persons
// @Produce json
// @Param id path int true "ID персоны"
// @Success 200 {object} dto.PersonDTO "Найденная персона"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} dto.ErrorResponse "Персона не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [get]
func (h *PersonHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, personHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, personHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	person, err := h.personService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toPersonDTO(person))
}

// GetByUserLogin возвращает персону по логину пользователя
// @Summary Получить персону по ID пользователя
// @Description Возвращает персону по идентификатору пользователя
// @Tags persons
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} dto.PersonDTO "Найденная персона"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат user_id"
// @Failure 404 {object} dto.ErrorResponse "Персона не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/user/{user_login} [get]
func (h *PersonHandler) GetByUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userLogin := r.PathValue("user_login")
	if userLogin == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, personHandlerCode, "user_login parameter missing"))
		return
	}

	person, err := h.personService.GetByUserLogin(ctx, userLogin)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toPersonDTO(person))
}

// Delete удаляет персону
// @Summary Удалить персону
// @Description Удаляет персону (мягкое или полное удаление)
// @Tags persons
// @Produce json
// @Param id path int true "ID персоны"
// @Param soft query bool false "Мягкое удаление (true/false)" default(true)
// @Success 204 "Персона успешно удалена"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID или параметра soft"
// @Failure 404 {object} dto.ErrorResponse "Персона не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [delete]
func (h *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, personHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, personHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	// Получаем параметр soft из query string
	softDelete := true
	softParam := r.URL.Query().Get("soft")
	if softParam != "" {
		soft, err := strconv.ParseBool(softParam)
		if err != nil {
			middleware.HandleError(w, r, appErr.NewLogicalError(err, personHandlerCode, "soft parameter must be boolean!"+err.Error()))
			return
		}
		softDelete = soft
	}

	person, err := h.personService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	err = h.personService.Delete(ctx, person, softDelete)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandler) toPersonDTO(person entities.Person) dto.PersonDTO {
	return dto.PersonDTO{
		ID:        person.ID,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,

		Firstname: person.Firstname,
		Lastname:  person.Lastname,
		Phone:     person.Phone,
		UserLogin: person.UserLogin,
	}
}
