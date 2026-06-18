package websocket

import (
	"strconv"

	fws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// wsConn ตั้งเป็น alias ของ connection จริงจาก gofiber/contrib/websocket
// (import ด้วยชื่อ fws เพื่อเลี่ยงชนกับชื่อ package ของเราเองที่ชื่อ websocket)
type wsConn = fws.Conn

// RegisterRoutes ผูก WebSocket endpoint เข้ากับ Fiber app
// เรียกจาก main หลัง SetupRoutes — ทำให้ ws เป็น delivery ที่ self-contained
// (package rest ไม่ต้องรู้จัก library websocket เลย)
func (h *Hub) RegisterRoutes(app *fiber.App) {
	// middleware: ผ่านเฉพาะ request ที่เป็น WebSocket upgrade จริง ๆ
	// ถ้าเปิดด้วย browser ปกติ (ไม่ใช่ ws) จะตอบ 426 Upgrade Required
	app.Use("/ws", func(c *fiber.Ctx) error {
		if fws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// ฟังข้อความสดของห้อง :id — ทุก connection จะถูกจัดการโดย hub
	app.Get("/ws/rooms/:id", fws.New(h.serve))
}

// serve ทำงานต่อ 1 connection: ลงทะเบียนเข้า hub แล้วเปิด read/write pump
func (h *Hub) serve(conn *wsConn) {
	roomID, err := strconv.Atoi(conn.Params("id"))
	if err != nil || roomID <= 0 {
		_ = conn.Close()
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		roomID: uint(roomID),
		send:   make(chan []byte, 256),
	}
	h.register <- client

	// writePump รันแยก goroutine: เป็น "ผู้เขียนคนเดียว" ลง connection (กันเขียนชนกัน)
	go client.writePump()

	// readPump รันใน goroutine ของ handler นี้: บล็อกอ่านจน connection ปิด
	// เราไม่ได้ใช้ข้อความที่ client พิมพ์เข้ามาทาง ws (การส่งจริงไปทาง REST/GraphQL)
	// อ่านไว้แค่เพื่อรับรู้ว่า connection ยังอยู่ไหม — พอปิดก็ unregister
	client.readPump()
}

// readPump อ่านจาก connection จนกว่าจะปิด/พัง แล้วถอนตัวออกจาก hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break // connection ปิดหรือ error → ออกจาก loop → unregister
		}
	}
}

// writePump ดึงข้อความจากคิว send แล้วเขียนลง connection ทีละอัน
// เมื่อ hub ปิดคิว (close(c.send)) loop จะจบเอง แล้วปิด connection
func (c *Client) writePump() {
	defer func() { _ = c.conn.Close() }()
	for payload := range c.send {
		if err := c.conn.WriteMessage(fws.TextMessage, payload); err != nil {
			break
		}
	}
}
