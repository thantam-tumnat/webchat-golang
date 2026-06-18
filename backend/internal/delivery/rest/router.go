package rest

import (
	"chatapp/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Handlers รวม handler ทุกตัวไว้ที่เดียว เพื่อส่งเข้า SetupRoutes
type Handlers struct {
	User    *UserHandler
	Room    *RoomHandler
	Message *MessageHandler
}

// NewHandlers ประกอบ handler ทั้งหมดจาก usecase (ทำใน main ตอน wiring)
func NewHandlers(
	userUC *usecases.UserUsecase,
	roomUC *usecases.RoomUsecase,
	messageUC *usecases.MessageUsecase,
) *Handlers {
	return &Handlers{
		User:    NewUserHandler(userUC),
		Room:    NewRoomHandler(roomUC),
		Message: NewMessageHandler(messageUC),
	}
}

// SetupRoutes ลงทะเบียน middleware + เส้นทาง API ทั้งหมด
// gqlHandler คือ endpoint ของ GraphQL (delivery layer อีกตัวที่ reuse usecase เดิม)
func SetupRoutes(app *fiber.App, h *Handlers, corsOrigins string, gqlHandler fiber.Handler) {
	// --- middleware ที่ทำงานก่อนถึง handler ทุก request ---
	app.Use(recover.New()) // ดัก panic ไม่ให้ server ตาย
	app.Use(logger.New())  // log ทุก request: method, path, status, latency
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins, // อนุญาตเฉพาะ origin ที่กำหนด (frontend)
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// health check สำหรับเช็คว่า server ยังทำงานอยู่
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// --- API routes (group /api) ---
	api := app.Group("/api")

	api.Post("/users", h.User.Create)

	api.Get("/rooms", h.Room.List)
	api.Post("/rooms", h.Room.Create)

	api.Get("/rooms/:id/messages", h.Message.List)
	api.Post("/rooms/:id/messages", h.Message.Send)

	// --- GraphQL endpoint ---
	// All: รองรับทั้ง GET (เปิด GraphiQL playground) และ POST (รัน query/mutation)
	// REST ด้านบนยังทำงานเหมือนเดิม GraphQL เป็นแค่ประตูเพิ่มอีกบาน
	app.All("/graphql", gqlHandler)
}
