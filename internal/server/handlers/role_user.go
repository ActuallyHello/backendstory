package handlers

// const (
// 	roleUserHandlerCode = "ROLE_USER_HANDLER"
// )

// type RoleUserHandler struct {
// 	ctx             context.Context
// 	validate        *validator.Validate
// 	roleUserService services.RoleUserService
// }

// func NewRoleUserHandler(
// 	ctx context.Context,
// 	roleUserService services.RoleUserService,
// ) *RoleUserHandler {
// 	return &RoleUserHandler{
// 		ctx:             ctx,
// 		validate:        validator.New(),
// 		roleUserService: roleUserService,
// 	}
// }

// // Create создает новую связь пользователя с ролью
// // @Summary Создать связь пользователя с ролью
// // @Description Создает новую связь между пользователем и ролью в системе
// // @Tags role-users
// // @Accept json
// // @Produce json
// // @Param request body dto.RoleUserCreateRequest true "Данные для создания связи"
// // @Success 201 {object} dto.RoleUserDTO "Созданная связь"
// // @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// // @Failure 409 {object} dto.ErrorResponse "Связь с такими roleID и userID уже существует"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /role-users [post]
// func (h *RoleUserHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	var req dto.RoleUserCreateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		middleware.HandleError(w, r, appErr.NewTechnicalError(err, roleUserHandlerCode, err.Error()))
// 		return
// 	}
// 	if err := h.validate.Struct(req); err != nil {
// 		details := common.CollectValidationDetails(err)
// 		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, err.Error()), details)
// 		return
// 	}

// 	roleUser := entities.RoleUser{
// 		RoleID: req.RoleID,
// 		UserID: req.UserID,
// 	}
// 	roleUser, err := h.roleUserService.Create(h.ctx, roleUser)
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(h.toRoleUserDTO(roleUser))
// }

// // GetByRoleID возвращает все связи по ID роли
// // @Summary Получить связи по ID роли
// // @Description Возвращает все связи пользователей с указанной ролью
// // @Tags role-users
// // @Produce json
// // @Param roleID path int true "ID роли"
// // @Success 200 {array} dto.RoleUserDTO "Найденные связи"
// // @Failure 400 {object} dto.ErrorResponse "Неверный формат roleID"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /role-users/role/{roleID} [get]
// func (h *RoleUserHandler) GetByRoleID(w http.ResponseWriter, r *http.Request) {
// 	reqRoleID := r.PathValue("roleID")
// 	if reqRoleID == "" {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleUserHandlerCode, "roleID parameter missing"))
// 		return
// 	}
// 	roleID, err := strconv.Atoi(reqRoleID)
// 	if err != nil {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, "roleID parameter must be integer!"+err.Error()))
// 		return
// 	}

// 	roleUsers, err := h.roleUserService.GetByRoleID(h.ctx, uint(roleID))
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	dtos := make([]dto.RoleUserDTO, len(roleUsers))
// 	for i, roleUser := range roleUsers {
// 		dtos[i] = h.toRoleUserDTO(roleUser)
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(dtos)
// }

// // GetByUserID возвращает все связи по ID пользователя
// // @Summary Получить связи по ID пользователя
// // @Description Возвращает все связи ролей с указанным пользователем
// // @Tags role-users
// // @Produce json
// // @Param userID path int true "ID пользователя"
// // @Success 200 {array} dto.RoleUserDTO "Найденные связи"
// // @Failure 400 {object} dto.ErrorResponse "Неверный формат userID"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /role-users/user/{userID} [get]
// func (h *RoleUserHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
// 	reqUserID := r.PathValue("userID")
// 	if reqUserID == "" {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleUserHandlerCode, "userID parameter missing"))
// 		return
// 	}
// 	userID, err := strconv.Atoi(reqUserID)
// 	if err != nil {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, "userID parameter must be integer!"+err.Error()))
// 		return
// 	}

// 	roleUsers, err := h.roleUserService.GetByUserID(h.ctx, uint(userID))
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	dtos := make([]dto.RoleUserDTO, len(roleUsers))
// 	for i, roleUser := range roleUsers {
// 		dtos[i] = h.toRoleUserDTO(roleUser)
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(dtos)
// }

// // GetByRoleIDAndUserID возвращает конкретную связь по roleID и userID
// // @Summary Получить связь по roleID и userID
// // @Description Возвращает конкретную связь между пользователем и ролью
// // @Tags role-users
// // @Produce json
// // @Param roleID path int true "ID роли"
// // @Param userID path int true "ID пользователя"
// // @Success 200 {object} dto.RoleUserDTO "Найденная связь"
// // @Failure 400 {object} dto.ErrorResponse "Неверный формат roleID или userID"
// // @Failure 404 {object} dto.ErrorResponse "Связь не найдена"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /role-users/role/{roleID}/user/{userID} [get]
// func (h *RoleUserHandler) GetByRoleIDAndUserID(w http.ResponseWriter, r *http.Request) {
// 	reqRoleID := r.PathValue("roleID")
// 	if reqRoleID == "" {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleUserHandlerCode, "roleID parameter missing"))
// 		return
// 	}
// 	reqUserID := r.PathValue("userID")
// 	if reqUserID == "" {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleUserHandlerCode, "userID parameter missing"))
// 		return
// 	}

// 	roleID, err := strconv.Atoi(reqRoleID)
// 	if err != nil {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, "roleID parameter must be integer!"+err.Error()))
// 		return
// 	}
// 	userID, err := strconv.Atoi(reqUserID)
// 	if err != nil {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, "userID parameter must be integer!"+err.Error()))
// 		return
// 	}

// 	roleUser, err := h.roleUserService.GetByRoleIDAndUserID(h.ctx, uint(roleID), uint(userID))
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(h.toRoleUserDTO(roleUser))
// }

// // Update обновляет связь пользователя с ролью
// // @Summary Обновить связь пользователя с ролью
// // @Description Обновляет существующую связь между пользователем и ролью
// // @Tags role-users
// // @Accept json
// // @Produce json
// // @Param request body dto.RoleUserUpdateRequest true "Данные для обновления"
// // @Success 204 "Связь успешно обновлена"
// // @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} dto.ErrorResponse "Связь не найдена"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /role-users [put]
// func (h *RoleUserHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	var req dto.RoleUserUpdateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		middleware.HandleError(w, r, appErr.NewTechnicalError(err, roleUserHandlerCode, err.Error()))
// 		return
// 	}
// 	if err := h.validate.Struct(req); err != nil {
// 		details := common.CollectValidationDetails(err)
// 		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, err.Error()), details)
// 		return
// 	}

// 	roleUser := entities.RoleUser{
// 		RoleID: req.RoleID,
// 		UserID: req.UserID,
// 	}
// 	_, err := h.roleUserService.Update(h.ctx, roleUser)
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// // Delete удаляет связь пользователя с ролью
// // @Summary Удалить связь пользователя с ролью
// // @Description Удаляет связь между пользователем и ролью
// // @Tags role-users
// // @Produce json
// // @Param roleID path int true "ID роли"
// // @Param userID path int true "ID пользователя"
// // @Success 204 "Связь успешно удалена"
// // @Failure 400 {object} dto.ErrorResponse "Неверный формат roleID или userID"
// // @Failure 404 {object} dto.ErrorResponse "Связь не найдена"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /role-users/role/{roleID}/user/{userID} [delete]
// func (h *RoleUserHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	reqRoleID := r.PathValue("roleID")
// 	if reqRoleID == "" {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleUserHandlerCode, "roleID parameter missing"))
// 		return
// 	}
// 	reqUserID := r.PathValue("userID")
// 	if reqUserID == "" {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleUserHandlerCode, "userID parameter missing"))
// 		return
// 	}

// 	roleID, err := strconv.Atoi(reqRoleID)
// 	if err != nil {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, "roleID parameter must be integer!"+err.Error()))
// 		return
// 	}
// 	userID, err := strconv.Atoi(reqUserID)
// 	if err != nil {
// 		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleUserHandlerCode, "userID parameter must be integer!"+err.Error()))
// 		return
// 	}

// 	roleUser, err := h.roleUserService.GetByRoleIDAndUserID(h.ctx, uint(roleID), uint(userID))
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}
// 	err = h.roleUserService.Delete(h.ctx, roleUser)
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// func (h *RoleUserHandler) toRoleUserDTO(roleUser entities.RoleUser) dto.RoleUserDTO {
// 	return dto.RoleUserDTO{
// 		ID:        roleUser.ID,
// 		CreatedAt: roleUser.CreatedAt,
// 		RoleID:    roleUser.RoleID,
// 		UserID:    roleUser.UserID,
// 	}
// }
