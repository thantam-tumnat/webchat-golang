package domain

import (
	"context"
	"time"
)

// Message คือ entity ของข้อความในห้องแชท
// มี field User แนบมาด้วย (preload) เพื่อให้ frontend แสดงชื่อผู้ส่งได้เลย
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoomID    uint      `gorm:"index;not null" json:"room_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// ความสัมพันธ์: 1 message เป็นของ 1 user
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// MessageRepository contract สำหรับเข้าถึงข้อความ
type MessageRepository interface {
	Create(ctx context.Context, msg *Message) error
	// FindByRoom คืน messages แบบแบ่งหน้า + จำนวนทั้งหมด (total)
	FindByRoom(ctx context.Context, roomID uint, page, limit int) ([]Message, int64, error)
}
