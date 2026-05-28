package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config เก็บค่า config ทั้งหมดของแอป โหลดมาจาก environment variables
// แยก config ออกจาก code ตามหลัก 12-factor app
type Config struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	CORSOrigins string
}

// Load อ่านค่าจากไฟล์ .env (ถ้ามี) แล้ว fallback ไปที่ environment variables ของระบบ
func Load() *Config {
	// โหลด .env ถ้าไม่มีก็ไม่เป็นไร (เช่นตอนรันบน production ใช้ env จริง)
	_ = godotenv.Load()

	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "chatapp"),
		DBPassword:  getEnv("DB_PASSWORD", "chatapp"),
		DBName:      getEnv("DB_NAME", "chatapp"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173"),
	}
}

// DSN สร้าง connection string สำหรับต่อ PostgreSQL ด้วย GORM
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
