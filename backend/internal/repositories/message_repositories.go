package repository

import (
	"context"

	"chatapp/internal/domain"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *domain.Message) error {
	if err := r.db.WithContext(ctx).Create(msg).Error; err != nil {
		return err
	}
	// preload user ของข้อความที่เพิ่งสร้าง เพื่อส่งกลับให้ frontend แสดงชื่อได้เลย
	return r.db.WithContext(ctx).Preload("User").First(msg, msg.ID).Error
}

// FindByRoom ดึงข้อความในห้องแบบแบ่งหน้า
// คืน (ข้อความหน้านั้น, จำนวนข้อความทั้งหมดในห้อง, error)
func (r *messageRepository) FindByRoom(ctx context.Context, roomID uint, page, limit int) ([]domain.Message, int64, error) {
	var messages []domain.Message
	var total int64

	// นับจำนวนทั้งหมดก่อน (ไว้คำนวณจำนวนหน้าฝั่ง frontend)
	if err := r.db.WithContext(ctx).
		Model(&domain.Message{}).
		Where("room_id = ?", roomID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("room_id = ?", roomID).
		Order("created_at DESC"). // ใหม่สุดก่อน แล้วค่อย reverse ฝั่ง frontend
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
