package repositories

import (
	"context"
	"errors"

	"chatapp/internal/entities"

	"gorm.io/gorm"
)

// userRepository คือ implementation จริงของ entities.UserRepository ที่ใช้ GORM
// เป็นชั้นเดียวที่ "รู้จัก" GORM — usecase ไม่รู้เรื่องนี้เลย
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository คืน implementation ในรูปของ interface (ทำให้ test/mock ได้ง่าย)
func NewUserRepository(db *gorm.DB) entities.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entities.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entities.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
