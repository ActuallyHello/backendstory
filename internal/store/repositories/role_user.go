package repositories

// type RoleUserRepository interface {
// 	Create(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error)
// 	Update(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error)
// 	Delete(ctx context.Context, roleUser entities.RoleUser) error

// 	FindByRoleID(ctx context.Context, roleID uint) ([]entities.RoleUser, error)
// 	FindByUserID(ctx context.Context, userID uint) ([]entities.RoleUser, error)
// 	FindByRoleIDAndUserID(ctx context.Context, roleID uint, userID uint) (entities.RoleUser, error)
// }

// type roleUserRepository struct {
// 	db *gorm.DB
// }

// func NewRoleUserRepository(db *gorm.DB) *roleUserRepository {
// 	return &roleUserRepository{db: db}
// }

// func (r *roleUserRepository) Create(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error) {
// 	if err := r.db.WithContext(ctx).Create(&roleUser).Error; err != nil {
// 		return entities.RoleUser{}, err
// 	}
// 	return roleUser, nil
// }

// func (r *roleUserRepository) Delete(ctx context.Context, roleUser entities.RoleUser) error {
// 	if err := r.db.WithContext(ctx).Delete(&roleUser).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *roleUserRepository) Update(ctx context.Context, roleUser entities.RoleUser) (entities.RoleUser, error) {
// 	if err := r.db.WithContext(ctx).Save(&roleUser).Error; err != nil {
// 		return entities.RoleUser{}, err
// 	}
// 	return roleUser, nil
// }

// func (r *roleUserRepository) FindByRoleID(ctx context.Context, roleID uint) ([]entities.RoleUser, error) {
// 	var roleUsers []entities.RoleUser
// 	if err := r.db.WithContext(ctx).Where("ROLEID = ?", roleID).Find(&roleUsers).Error; err != nil {
// 		return nil, err
// 	}
// 	return roleUsers, nil
// }

// func (r *roleUserRepository) FindByUserID(ctx context.Context, userID uint) ([]entities.RoleUser, error) {
// 	var roleUsers []entities.RoleUser
// 	if err := r.db.WithContext(ctx).Where("USERID = ?", userID).Find(&roleUsers).Error; err != nil {
// 		return nil, err
// 	}
// 	return roleUsers, nil
// }

// func (r *roleUserRepository) FindByRoleIDAndUserID(ctx context.Context, roleID uint, userID uint) (entities.RoleUser, error) {
// 	var roleUser entities.RoleUser
// 	if err := r.db.WithContext(ctx).Where("ROLEID = ? AND USERID = ?", roleID, userID).First(&roleUser).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return entities.RoleUser{}, common.NewNotFoundError("role user not found")
// 		}
// 		return entities.RoleUser{}, err
// 	}
// 	return roleUser, nil
// }
