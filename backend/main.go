package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"chatapp/internal/config"
	"chatapp/internal/delivery/graphql"
	"chatapp/internal/delivery/rest"
	"chatapp/internal/delivery/websocket"
	"chatapp/internal/infrastructure/database"
	"chatapp/internal/repositories"
	"chatapp/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1) โหลด config จาก .env / environment
	cfg := config.Load()

	// 2) ต่อ database
	db, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("❌ ต่อ database ไม่ได้: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("❌ migrate ไม่สำเร็จ: %v", err)
	}

	// 3) Dependency Injection — สร้างจากชั้นล่างขึ้นบน
	//    repository (คุย DB)  ->  usecase (business logic)  ->  handler (HTTP)
	userRepo := repositories.NewUserRepository(db)
	roomRepo := repositories.NewRoomRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// WebSocket hub — เริ่มก่อน usecase เพราะ messageUC ต้องใช้เป็น notifier
	// รัน Run() ใน goroutine เดียวคอยจัดการ connection + broadcast ตลอดอายุ server
	hub := websocket.NewHub()
	go hub.Run()

	userUC := usecases.NewUserUsecase(userRepo)
	roomUC := usecases.NewRoomUsecase(roomRepo)
	messageUC := usecases.NewMessageUsecase(messageRepo, roomRepo, userRepo, hub)

	handlers := rest.NewHandlers(userUC, roomUC, messageUC)

	// GraphQL: ใช้ usecase ชุดเดียวกับ REST แล้วประกอบเป็น schema + handler
	gqlResolver := graphql.NewResolver(userUC, roomUC, messageUC)
	gqlSchema, err := graphql.NewSchema(gqlResolver)
	if err != nil {
		log.Fatalf("❌ สร้าง GraphQL schema ไม่สำเร็จ: %v", err)
	}
	gqlHandler := graphql.NewHandler(gqlSchema)

	// 4) สร้าง Fiber app พร้อม custom error handler กลาง
	app := fiber.New(fiber.Config{
		ErrorHandler: rest.ErrorHandler,
	})
	rest.SetupRoutes(app, handlers, cfg.CORSOrigins, gqlHandler)

	// WebSocket route: GET /ws/rooms/:id — ฟังข้อความสดของห้องนั้น
	// ลงทะเบียนหลัง SetupRoutes เพื่อให้ได้ middleware กลาง (recover/logger/cors) ด้วย
	hub.RegisterRoutes(app)

	// 5) รัน server ใน goroutine เพื่อให้ main รอ signal ปิดได้
	go func() {
		addr := ":" + cfg.AppPort
		log.Printf("🚀 server กำลังรันที่ http://localhost%s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("❌ server หยุดทำงาน: %v", err)
		}
	}()

	// 6) Graceful shutdown — รอ signal (Ctrl+C / docker stop)
	//    แล้วปิด server โดยรอ request ที่ค้างอยู่ให้เสร็จก่อน
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 กำลังปิด server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("error ตอนปิด server: %v", err)
	}
	log.Println("✅ ปิด server เรียบร้อย")
}
