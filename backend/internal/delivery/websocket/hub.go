// Package websocket เป็น delivery layer อีกบาน (เทียร์เดียวกับ rest, graphql)
// ต่างกันตรงรูปแบบ: REST/GraphQL เป็น "ถาม-ตอบ" ส่วนนี่เป็นท่อค้างไว้ที่ server "ดัน" ข้อมูลมาหา client ได้เอง
// ใช้เฉพาะ usecase ที่ข้อมูลเปลี่ยนแบบ real-time → ในโปรเจคนี้คือ "ข้อความใหม่ในห้องแชท"
package websocket

import (
	"encoding/json"
	"log"

	"chatapp/internal/entities"
)

// Client = หนึ่ง connection ของผู้ใช้ที่กำลังเปิดอยู่ในห้องหนึ่ง
type Client struct {
	hub    *Hub
	conn   *wsConn      // connection จริง (websocket)
	roomID uint         // อยู่ห้องไหน
	send   chan []byte  // คิวข้อความที่รอเขียนออกไปยัง connection นี้
}

// Hub = ตัวกลางจัดการทุก connection แยกตามห้อง + กระจายข้อความ (broadcast)
//
// แก้ปัญหา data race ด้วย "goroutine เดียว" (Run) เป็นเจ้าของ map rooms
// ทุกการแก้ไข map วิ่งผ่าน channel เข้ามาทีละ event → ไม่ต้องใช้ mutex เอง
type Hub struct {
	rooms      map[uint]map[*Client]bool // ห้อง -> เซ็ตของ client ที่ออนไลน์ในห้องนั้น
	register   chan *Client              // มีคนเข้าห้อง
	unregister chan *Client              // มีคนออก/หลุด
	broadcast  chan broadcast            // มีข้อความใหม่ต้องกระจาย
}

type broadcast struct {
	roomID  uint
	payload []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan broadcast, 256),
	}
}

// Run เป็น loop หลักของ hub — รัน 1 goroutine ตลอดอายุ server
// รับ event ทีละตัวจาก channel จึงเข้าถึง map ได้อย่างปลอดภัยโดยไม่มี race
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			if h.rooms[c.roomID] == nil {
				h.rooms[c.roomID] = make(map[*Client]bool)
			}
			h.rooms[c.roomID][c] = true

		case c := <-h.unregister:
			h.removeClient(c)

		case b := <-h.broadcast:
			for c := range h.rooms[b.roomID] {
				select {
				case c.send <- b.payload: // ส่งเข้าคิวของ client
				default:
					// client รับไม่ทัน (ช้า/ค้าง) → ตัดทิ้ง กัน hub โดนบล็อก
					h.removeClient(c)
				}
			}
		}
	}
}

// removeClient เอา client ออกจากห้อง + ปิดคิว send (ต้องเรียกจาก goroutine ของ Run เท่านั้น)
func (h *Hub) removeClient(c *Client) {
	if clients, ok := h.rooms[c.roomID]; ok {
		if _, ok := clients[c]; ok {
			delete(clients, c)
			close(c.send)
			if len(clients) == 0 {
				delete(h.rooms, c.roomID) // ห้องว่างแล้วเก็บกวาด
			}
		}
	}
}

// NotifyMessage = implement entities.MessageNotifier
// usecase เรียกตัวนี้หลังบันทึกข้อความสำเร็จ โดยไม่รู้เลยว่าเบื้องหลังเป็น WebSocket
func (h *Hub) NotifyMessage(roomID uint, msg *entities.Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws: marshal message ไม่สำเร็จ: %v", err)
		return
	}
	h.broadcast <- broadcast{roomID: roomID, payload: payload}
}
