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

// Create создает новую персону
// @Summary Создать персону
// @Description Создает новую персону в системе
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePersonRequest true "Данные для создания персоны"
// @Success 201 {object} dto.PersonDTO "Созданная персона"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Персона с таким user_login уже существует"
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
	json.NewEncoder(w).Encode(dto.ToPersonDTO(person))
}

// GetAll возвращает всех персон
// @Summary Получить всех персон
// @Description Возвращает список всех персон в системе
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PersonDTO "Список персон"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
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
		dtos = append(dtos, dto.ToPersonDTO(person))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает персону по ID
// @Summary Получить персону по ID
// @Description Возвращает персону по указанному идентификатору
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID персоны"
// @Success 200 {object} dto.PersonDTO "Персона"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
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
	json.NewEncoder(w).Encode(dto.ToPersonDTO(person))
}

// GetByUserLogin возвращает персону по логину пользователя
// @Summary Получить персону по логину пользователя
// @Description Возвращает персону по указанному логину пользователя
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_login path string true "Логин пользователя"
// @Success 200 {object} dto.PersonDTO "Персона"
// @Failure 400 {object} dto.ErrorResponse "Неверный логин пользователя"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
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
	json.NewEncoder(w).Encode(dto.ToPersonDTO(person))
}

// Delete удаляет персону
// @Summary Удалить персону
// @Description Удаляет персону по указанному идентификатору. Поддерживает мягкое удаление (soft delete)
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID персоны"
// @Param soft query boolean false "Флаг мягкого удаления (true/false)" default(true)
// @Success 204 "Успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID или параметр soft"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
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

// GetWithSearchCriteria выполняет поиск людей по критериям
// @Summary Поиск людей
// @Description Выполняет поиск людей по заданным критериям
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.PersonDTO "Список найденных людей"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/search [post]
func (h *PersonHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, personHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, personHandlerCode, err.Error()), details)
		return
	}

	persons, err := h.personService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.PersonDTO, 0, len(persons))
	for _, person := range persons {
		dtos = append(dtos, dto.ToPersonDTO(person))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
