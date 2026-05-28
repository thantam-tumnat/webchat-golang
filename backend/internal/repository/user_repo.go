package repository

import (
	"context"
	"errors"

	"chatapp/internal/domain"

	"gorm.io/gorm"
)

// userRepository คือ implementation จริงของ domain.UserRepository ที่ใช้ GORM
// เป็นชั้นเดียวที่ "รู้จัก" GORM — usecase ไม่รู้เรื่องนี้เลย
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository คืน implementation ในรูปของ interface (ทำให้ test/mock ได้ง่าย)
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
