package repositories

// type UserRepository interface {
// 	Create(ctx context.Context, user entities.User) (entities.User, error)
// 	Update(ctx context.Context, user entities.User) (entities.User, error)
// 	Delete(ctx context.Context, user entities.User, soft bool) error

// 	FindById(ctx context.Context, id uint) (entities.User, error)
// 	FindByEmail(ctx context.Context, email string) (entities.User, error)
// }

// type userRepository struct {
// 	db *gorm.DB
// }

// func NewUserRepository(db *gorm.DB) *userRepository {
// 	return &userRepository{db: db}
// }

// // Create создает нового пользователя
// func (r *userRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
// 	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
// 		return entities.User{}, err
// 	}
// 	return user, nil
// }

// // Delete выполняет удаление пользователя
// func (r *userRepository) Delete(ctx context.Context, user entities.User, soft bool) error {
// 	if soft {
// 		user.DeletedAt = sql.NullTime{
// 			Time:  time.Now(),
// 			Valid: true,
// 		}
// 		_, err := r.Update(ctx, user)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	} else {
// 		if err := r.db.WithContext(ctx).Delete(&user).Error; err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// }

// // Update обновляет данные пользователя
// func (r *userRepository) Update(ctx context.Context, user entities.User) (entities.User, error) {
// 	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
// 		return entities.User{}, err
// 	}
// 	return user, nil
// }

// // FindById ищет пользователя по ID
// func (r *userRepository) FindById(ctx context.Context, id uint) (entities.User, error) {
// 	var user entities.User
// 	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return entities.User{}, common.NewNotFoundError("user not found")
// 		}
// 		return entities.User{}, err
// 	}
// 	return user, nil
// }

// // FindByEmail ищет пользователя по email
// func (r *userRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
// 	var user entities.User
// 	if err := r.db.WithContext(ctx).Where("EMAIL = ?", email).First(&user).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return entities.User{}, common.NewNotFoundError("user not found")
// 		}
// 		return entities.User{}, err
// 	}
// 	return user, nil
// }
