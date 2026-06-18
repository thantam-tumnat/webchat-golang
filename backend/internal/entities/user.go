package entities

import (
	"context"
	"time"
)

// User คือ entity ของผู้ใช้ (Phase 1 มีแค่ username พอ ยังไม่มี password)
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository เป็น interface (contract) ของการเข้าถึงข้อมูล user
// usecase จะรู้จักแค่ interface นี้ ไม่รู้ว่าเบื้องหลังใช้ GORM/Postgres หรืออะไร
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}
