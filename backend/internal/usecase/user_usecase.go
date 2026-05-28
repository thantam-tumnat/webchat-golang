package usecase

import (
	"context"
	"errors"

	"chatapp/internal/domain"
)

// UserUsecase เก็บ business logic เกี่ยวกับ user
type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// CreateOrGet: Phase 1 ยังไม่มีระบบ login จริง
// ถ้า username นี้มีอยู่แล้วก็คืน user เดิม (idempotent) ถ้ายังไม่มีก็สร้างใหม่
// ทำให้ frontend แค่ส่ง username มาก็ใช้งานต่อได้เลย
func (uc *UserUsecase) CreateOrGet(ctx context.Context, username string) (*domain.User, error) {
	existing, err := uc.userRepo.FindByUsername(ctx, username)
	if err == nil {
		return existing, nil
	}
	// ถ้า error เป็นอย่างอื่นที่ไม่ใช่ "ไม่เจอ" ให้ส่ง error ออกไป
	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != domain.ErrUserNotFound.Code {
		return nil, err
	}

	user := &domain.User{Username: username}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
