package usecases

import (
	"context"

	"chatapp/internal/entities"
)

type MessageUsecase struct {
	messageRepo entities.MessageRepository
	roomRepo    entities.RoomRepository
	userRepo    entities.UserRepository
}

func NewMessageUsecase(
	messageRepo entities.MessageRepository,
	roomRepo entities.RoomRepository,
	userRepo entities.UserRepository,
) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
	}
}

// List ดึงข้อความในห้อง (ตรวจก่อนว่าห้องมีอยู่จริง)
func (uc *MessageUsecase) List(ctx context.Context, roomID uint, page, limit int) ([]entities.Message, int64, error) {
	if _, err := uc.roomRepo.FindByID(ctx, roomID); err != nil {
		return nil, 0, err
	}
	return uc.messageRepo.FindByRoom(ctx, roomID, page, limit)
}

// Send ส่งข้อความใหม่ (ตรวจว่าห้อง + ผู้ส่ง มีอยู่จริง)
func (uc *MessageUsecase) Send(ctx context.Context, roomID, userID uint, content string) (*entities.Message, error) {
	if _, err := uc.roomRepo.FindByID(ctx, roomID); err != nil {
		return nil, err
	}
	if _, err := uc.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}

	msg := &entities.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	}
	if err := uc.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
