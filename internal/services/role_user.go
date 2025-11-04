package services

// const (
// 	roleUserServiceCode = "ROLE_USER_SERVICE"
// )

// type RoleUserService interface {
// 	Create(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error)
// 	Update(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error)
// 	Delete(ctx context.Context, roleUser entities.RoleUser) error

// 	GetByRoleID(ctx context.Context, roleID uint) ([]entities.RoleUser, error)
// 	GetByUserID(ctx context.Context, userID uint) ([]entities.RoleUser, error)
// 	GetByRoleIDAndUserID(ctx context.Context, roleID uint, userID uint) (entities.RoleUser, error)
// }

// type roleUserService struct {
// 	roleUserRepo repositories.RoleUserRepository
// }

// func NewRoleUserService(
// 	roleUserRepo repositories.RoleUserRepository,
// ) *roleUserService {
// 	return &roleUserService{
// 		roleUserRepo: roleUserRepo,
// 	}
// }

// func (s *roleUserService) Create(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error) {
// 	// Проверка существования связи с такими RoleID и UserID
// 	existing, err := s.GetByRoleIDAndUserID(ctx, roleUser.RoleID, roleUser.UserID)
// 	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
// 		return entities.RoleUser{}, err
// 	}
// 	if existing.ID > 0 {
// 		slog.Error("RoleUser already exists!", "error", err, "roleID", roleUser.RoleID, "userID", roleUser.UserID)
// 		return entities.RoleUser{}, appError.NewLogicalError(nil, roleUserServiceCode, fmt.Sprintf("RoleUser with roleID = %d and userID = %d already exists!", roleUser.RoleID, roleUser.UserID))
// 	}

// 	// Создаем запись
// 	created, err := s.roleUserRepo.Create(ctx, roleUser)
// 	if err != nil {
// 		slog.Error("Create roleUser failed", "error", err, "roleID", roleUser.RoleID, "userID", roleUser.UserID)
// 		return entities.RoleUser{}, appError.NewTechnicalError(err, roleUserServiceCode, err.Error())
// 	}
// 	slog.Info("RoleUser created", "roleID", created.RoleID, "userID", created.UserID)
// 	return created, nil
// }

// func (s *roleUserService) Update(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error) {
// 	existing, err := s.GetByRoleIDAndUserID(ctx, roleUser.RoleID, roleUser.UserID)
// 	if err != nil {
// 		return entities.RoleUser{}, err
// 	}

// 	updated, err := s.roleUserRepo.Update(ctx, existing)
// 	if err != nil {
// 		slog.Error("Update roleUser failed", "error", err, "roleID", roleUser.RoleID, "userID", roleUser.UserID)
// 		return entities.RoleUser{}, appError.NewTechnicalError(err, roleUserServiceCode, err.Error())
// 	}
// 	return updated, nil
// }

// func (s *roleUserService) Delete(ctx context.Context, roleUser entities.RoleUser) error {
// 	err := s.roleUserRepo.Delete(ctx, roleUser)
// 	if err != nil {
// 		slog.Error("Failed to delete roleUser", "error", err, "roleID", roleUser.RoleID, "userID", roleUser.UserID)
// 		return appError.NewTechnicalError(err, roleUserServiceCode, err.Error())
// 	}
// 	slog.Info("Deleted roleUser", "roleID", roleUser.RoleID, "userID", roleUser.UserID)
// 	return nil
// }

// func (s *roleUserService) GetByRoleID(ctx context.Context, roleID uint) ([]entities.RoleUser, error) {
// 	roleUsers, err := s.roleUserRepo.FindByRoleID(ctx, roleID)
// 	if err != nil {
// 		slog.Error("Failed to find roleUsers by RoleID", "error", err, "roleID", roleID)
// 		return nil, appError.NewTechnicalError(err, roleUserServiceCode, err.Error())
// 	}
// 	return roleUsers, nil
// }

// func (s *roleUserService) GetByUserID(ctx context.Context, userID uint) ([]entities.RoleUser, error) {
// 	roleUsers, err := s.roleUserRepo.FindByUserID(ctx, userID)
// 	if err != nil {
// 		slog.Error("Failed to find roleUsers by UserID", "error", err, "userID", userID)
// 		return nil, appError.NewTechnicalError(err, roleUserServiceCode, err.Error())
// 	}
// 	return roleUsers, nil
// }

// func (s *roleUserService) GetByRoleIDAndUserID(ctx context.Context, roleID uint, userID uint) (entities.RoleUser, error) {
// 	roleUser, err := s.roleUserRepo.FindByRoleIDAndUserID(ctx, roleID, userID)
// 	if err != nil {
// 		slog.Error("Failed to find roleUser by RoleID and UserID", "error", err, "roleID", roleID, "userID", userID)
// 		if errors.Is(err, &common.NotFoundError{}) {
// 			return entities.RoleUser{}, appError.NewLogicalError(err, roleUserServiceCode, err.Error())
// 		}
// 		return entities.RoleUser{}, appError.NewTechnicalError(err, roleUserServiceCode, err.Error())
// 	}
// 	return roleUser, nil
// }
