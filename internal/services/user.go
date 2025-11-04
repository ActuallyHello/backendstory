package services

// const (
// 	userServiceCode = "USER_SERVICE"
// )

// type UserService interface {
// 	Create(ctx context.Context, user entities.User) (entities.User, error)
// 	Update(ctx context.Context, user entities.User) (entities.User, error)
// 	Delete(ctx context.Context, user entities.User, soft bool) error

// 	FindById(ctx context.Context, id uint) (entities.User, error)
// 	FindByEmail(ctx context.Context, email string) (entities.User, error)
// }

// type userService struct {
// 	userRepo repositories.UserRepository
// }

// func NewUserService(
// 	userRepo repositories.UserRepository,
// ) *userService {
// 	return &userService{
// 		userRepo: userRepo,
// 	}
// }

// // Create создает нового пользователя
// func (s *userService) Create(ctx context.Context, user entities.User) (entities.User, error) {
// 	// Проверка существования пользователя с таким email
// 	existing, err := s.FindByEmail(ctx, user.Email)
// 	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
// 		return entities.User{}, err
// 	}
// 	if existing.ID > 0 {
// 		slog.Error("User with this email already exists!", "error", err, "email", user.Email)
// 		return entities.User{}, appError.NewLogicalError(nil, userServiceCode, fmt.Sprintf("User with email = %s already exists!", user.Email))
// 	}

// 	// Создаем запись
// 	created, err := s.userRepo.Create(ctx, user)
// 	if err != nil {
// 		slog.Error("Create user failed", "error", err, "email", user.Email)
// 		return entities.User{}, appError.NewTechnicalError(err, userServiceCode, err.Error())
// 	}
// 	slog.Info("User created", "email", created.Email, "user_id", created.ID)
// 	return created, nil
// }

// // Update обновляет существующего пользователя
// func (s *userService) Update(ctx context.Context, user entities.User) (entities.User, error) {
// 	// Проверяем существование пользователя
// 	existing, err := s.FindById(ctx, user.ID)
// 	if err != nil {
// 		return entities.User{}, err
// 	}

// 	// Проверяем уникальность email если он изменился
// 	if existing.Email != user.Email {
// 		existingByEmail, err := s.FindByEmail(ctx, user.Email)
// 		if err != nil && errors.Is(err, &appError.TechnicalError{}) {
// 			return entities.User{}, err
// 		}
// 		if existingByEmail.ID > 0 {
// 			slog.Error("User with this email already exists!", "error", err, "email", user.Email)
// 			return entities.User{}, appError.NewLogicalError(nil, userServiceCode, fmt.Sprintf("User with email = %s already exists!", user.Email))
// 		}
// 	}

// 	updated, err := s.userRepo.Update(ctx, user)
// 	if err != nil {
// 		slog.Error("Update user failed", "error", err, "user_id", user.ID, "email", user.Email)
// 		return entities.User{}, appError.NewTechnicalError(err, userServiceCode, err.Error())
// 	}
// 	slog.Info("User updated", "user_id", updated.ID, "email", updated.Email)
// 	return updated, nil
// }

// // Delete удаляет пользователя (мягко или полностью)
// func (s *userService) Delete(ctx context.Context, user entities.User, soft bool) error {
// 	err := s.userRepo.Delete(ctx, user, soft)
// 	if err != nil {
// 		slog.Error("Failed to delete user", "error", err, "user_id", user.ID, "soft", soft)
// 		return appError.NewTechnicalError(err, userServiceCode, err.Error())
// 	}
// 	slog.Info("Deleted user", "user_id", user.ID, "soft", soft)
// 	return nil
// }

// // FindById ищет пользователя по ID
// func (s *userService) FindById(ctx context.Context, id uint) (entities.User, error) {
// 	user, err := s.userRepo.FindById(ctx, id)
// 	if err != nil {
// 		slog.Error("Failed to find user by ID", "error", err, "id", id)
// 		if errors.Is(err, &common.NotFoundError{}) {
// 			return entities.User{}, appError.NewLogicalError(err, userServiceCode, err.Error())
// 		}
// 		return entities.User{}, appError.NewTechnicalError(err, userServiceCode, err.Error())
// 	}
// 	return user, nil
// }

// // FindByEmail ищет пользователя по email
// func (s *userService) FindByEmail(ctx context.Context, email string) (entities.User, error) {
// 	user, err := s.userRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		slog.Error("Failed to find user by email", "error", err, "email", email)
// 		if errors.Is(err, &common.NotFoundError{}) {
// 			return entities.User{}, appError.NewLogicalError(err, userServiceCode, err.Error())
// 		}
// 		return entities.User{}, appError.NewTechnicalError(err, userServiceCode, err.Error())
// 	}
// 	return user, nil
// }
