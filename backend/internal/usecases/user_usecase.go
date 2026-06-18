package usecases

import (
	"context"
	"errors"

	"chatapp/internal/entities"
)

// UserUsecase เก็บ business logic เกี่ยวกับ user
type UserUsecase struct {
	userRepo entities.UserRepository
}

func NewUserUsecase(userRepo entities.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// CreateOrGet: Phase 1 ยังไม่มีระบบ login จริง
// ถ้า username นี้มีอยู่แล้วก็คืน user เดิม (idempotent) ถ้ายังไม่มีก็สร้างใหม่
// ทำให้ frontend แค่ส่ง username มาก็ใช้งานต่อได้เลย
func (uc *UserUsecase) CreateOrGet(ctx context.Context, username string) (*entities.User, error) {
	existing, err := uc.userRepo.FindByUsername(ctx, username)
	if err == nil {
		return existing, nil
	}
	// ถ้า error เป็นอย่างอื่นที่ไม่ใช่ "ไม่เจอ" ให้ส่ง error ออกไป
	var appErr *entities.AppError
	if !errors.As(err, &appErr) || appErr.Code != entities.ErrUserNotFound.Code {
		return nil, err
	}

	user := &entities.User{Username: username}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
