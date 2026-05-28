package database

import (
	"log"
	"time"

	"chatapp/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgres เปิด connection ไป PostgreSQL ด้วย GORM
func NewPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // log เฉพาะ query ที่ช้า/error
	})
	if err != nil {
		return nil, err
	}

	// ตั้งค่า connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// AutoMigrate สร้าง/อัปเดต table ตาม struct ใน domain ให้อัตโนมัติ
// (Phase 1 ใช้วิธีนี้เพื่อให้รันได้ทันที — Phase หลังจะเปลี่ยนไปใช้ golang-migrate)
func AutoMigrate(db *gorm.DB) error {
	log.Println("🔧 running auto-migration...")
	return db.AutoMigrate(
		&domain.User{},
		&domain.Room{},
		&domain.Message{},
	)
}
