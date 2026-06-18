package usecases

import (
	"context"

	"chatapp/internal/entities"
)

type MessageUsecase struct {
	messageRepo entities.MessageRepository
	roomRepo    entities.RoomRepository
	userRepo    entities.UserRepository
	notifier    entities.MessageNotifier // ช่องทาง broadcast real-time (อาจเป็น nil ได้)
}

func NewMessageUsecase(
	messageRepo entities.MessageRepository,
	roomRepo entities.RoomRepository,
	userRepo entities.UserRepository,
	notifier entities.MessageNotifier,
) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
		notifier:    notifier,
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
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
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

	// แนบ user ที่เพิ่งตรวจไปแล้ว เพื่อให้ฝั่งรับแสดงชื่อผู้ส่งได้ทันที (เหมือนตอน List ที่ preload มา)
	msg.User = user

	// กระจายข้อความสด ๆ ให้ทุกคนในห้อง ถ้ามี notifier (WebSocket)
	// ส่งผ่าน REST หรือ GraphQL ก็ broadcast เหมือนกัน เพราะวิ่งผ่าน usecase เดียวกันนี้
	if uc.notifier != nil {
		uc.notifier.NotifyMessage(roomID, msg)
	}
	return msg, nil
}
