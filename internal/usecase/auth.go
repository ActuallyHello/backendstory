package usecase

// const (
// 	authUsecaseCode  = "AUTH_USECASE"
// 	enumStatusCode   = "Status"
// 	activeStatusCode = "Active"
// )

// type AuthUsecase interface {
// 	RegisterUser(ctx context.Context, userDTO dto.RegisterUserRequest) (dto.UserDTO, error)
// 	AuthtorizeUser(ctx context.Context, loginRequest dto.LoginRequest) (dto.UserDTO, error)
// }

// type authUsecase struct {
// 	ctx                     context.Context
// 	userService             services.UserService
// 	enumerationService      services.EnumService
// 	enumerationValueService services.EnumValueService
// }

// func NewAuthUsecase(
// 	ctx context.Context,
// 	enumerationService services.EnumService,
// 	enumerationValueService services.EnumValueService,
// 	userService services.UserService,
// ) *authUsecase {
// 	return &authUsecase{
// 		ctx:                     ctx,
// 		userService:             userService,
// 		enumerationService:      enumerationService,
// 		enumerationValueService: enumerationValueService,
// 	}
// }

// func (u *authUsecase) RegisterUser(ctx context.Context, userDTO dto.RegisterUserRequest) (dto.UserDTO, error) {
// 	hashedPassword, err := auth.GenerateHash(userDTO.Password)
// 	if err != nil {
// 		slog.Error("Failed to hash password!", "error", err)
// 		return dto.UserDTO{}, appErr.NewTechnicalError(err, authUsecaseCode, err.Error())
// 	}

// 	statusEnum, err := u.enumerationService.GetByCode(u.ctx, enumStatusCode)
// 	if err != nil {
// 		return dto.UserDTO{}, err
// 	}
// 	activeEnumVal, err := u.enumerationValueService.GetByCodeAndEnumID(u.ctx, activeStatusCode, statusEnum.ID)
// 	if err != nil {
// 		return dto.UserDTO{}, err
// 	}

// 	user := entities.User{
// 		Email:    userDTO.Email,
// 		Password: hashedPassword,
// 		StatusID: activeEnumVal.ID,
// 	}
// 	user, err = u.userService.Create(u.ctx, user)
// 	if err != nil {
// 		return dto.UserDTO{}, err
// 	}

// 	return dto.UserDTO{
// 		ID:        user.ID,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 		Email:     user.Email,
// 		StatusID:  user.StatusID,
// 	}, nil
// }

// func (u *authUsecase) AuthtorizeUser(ctx context.Context, loginRequest dto.LoginRequest) (dto.UserDTO, error) {
// 	user, err := u.userService.FindByEmail(u.ctx, loginRequest.Email)
// 	if err != nil {
// 		return dto.UserDTO{}, err
// 	}
// 	isValid, err := auth.VerifyPassword(loginRequest.Password, user.Password)
// 	if err != nil {
// 		return dto.UserDTO{}, err
// 	}
// 	if !isValid {
// 		slog.Error("Invalid passowrd attempt", "email", loginRequest.Email)
// 		return dto.UserDTO{}, appErr.NewLogicalError(nil, authUsecaseCode, fmt.Sprintf("Invalid password for email=%s!", loginRequest.Email))
// 	}
// 	return dto.UserDTO{
// 		ID:        user.ID,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 		Email:     user.Email,
// 		StatusID:  user.StatusID,
// 	}, nil
// }
