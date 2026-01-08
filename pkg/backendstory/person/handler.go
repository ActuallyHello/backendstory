package person

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	personHandlerCode = "PERSON_HANDLER"
)

type PersonHandler struct {
	validate      *validator.Validate
	personService PersonService
}

func NewPersonHandler(
	personService PersonService,
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
// @Param request body CreatePersonRequest true "Данные для создания персоны"
// @Success 201 {object} PersonDTO "Созданная персона"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Персона с таким user_login уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons [post]
// @Id createPerson
func (h *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, personHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, personHandlerCode, err.Error()), details)
		return
	}

	person := Person{
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Phone:     req.Phone,
		UserLogin: req.UserLogin,
	}
	person, err := h.personService.Create(ctx, person)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToPersonDTO(person))
}

// GetAll возвращает всех персон
// @Summary Получить всех персон
// @Description Возвращает список всех персон в системе
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} PersonDTO "Список персон"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons [get]
// @Id getPersonAll
func (h *PersonHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	persons, err := h.personService.GetAll(ctx)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]PersonDTO, 0, len(persons))
	for _, person := range persons {
		dtos = append(dtos, ToPersonDTO(person))
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
// @Success 200 {object} PersonDTO "Персона"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Персона не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [get]
// @Id getPersonById
func (h *PersonHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, personHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, personHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	person, err := h.personService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToPersonDTO(person))
}

// GetByUserLogin возвращает персону по логину пользователя
// @Summary Получить персону по логину пользователя
// @Description Возвращает персону по указанному логину пользователя
// @Tags Persons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_login path string true "Логин пользователя"
// @Success 200 {object} PersonDTO "Персона"
// @Failure 400 {object} core.ErrorResponse "Неверный логин пользователя"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Персона не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/user/{user_login} [get]
// @Id getPersonByUserLogin
func (h *PersonHandler) GetByUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userLogin := r.PathValue("user_login")
	if userLogin == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, personHandlerCode, "Отсуствует логин пользователя"))
		return
	}

	person, err := h.personService.GetByUserLogin(ctx, userLogin)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToPersonDTO(person))
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
// @Failure 400 {object} core.ErrorResponse "Неверный ID или параметр soft"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Персона не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [delete]
// @Id deletePerson
func (h *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, personHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, personHandlerCode, "ИД параметр должен быть числовым! "+err.Error()))
		return
	}

	// Получаем параметр soft из query string
	softDelete := true
	softParam := r.URL.Query().Get("soft")
	if softParam != "" {
		soft, err := strconv.ParseBool(softParam)
		if err != nil {
			core.HandleError(w, r, core.NewLogicalError(err, personHandlerCode, "Признак мягкого удаления должен быть булевым!"+err.Error()))
			return
		}
		softDelete = soft
	}

	person, err := h.personService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	err = h.personService.Delete(ctx, person, softDelete)
	if err != nil {
		core.HandleError(w, r, err)
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
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} PersonDTO "Список найденных людей"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/search [post]
// @Id searchPerson
func (h *PersonHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, personHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, personHandlerCode, err.Error()), details)
		return
	}

	persons, err := h.personService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]PersonDTO, 0, len(persons))
	for _, person := range persons {
		dtos = append(dtos, ToPersonDTO(person))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
