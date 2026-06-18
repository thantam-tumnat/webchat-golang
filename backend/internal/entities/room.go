package entities

import (
	"context"
	"time"
)

// Room คือ entity ของห้องแชท
type Room struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RoomRepository contract สำหรับเข้าถึงข้อมูลห้อง
type RoomRepository interface {
	Create(ctx context.Context, room *Room) error
	FindAll(ctx context.Context) ([]Room, error)
	FindByID(ctx context.Context, id uint) (*Room, error)
}
