# 💬 Real-time Chat App (Go + React)

A real-time group chat application built with Go, React, REST APIs, GraphQL, and WebSockets.

จุดประสงค์ของโปรเจคนี้ไม่ใช่แค่ทำแชทให้ส่งข้อความได้ แต่เป็นการทดลองว่า REST, GraphQL และ WebSocket สามารถอยู่ร่วมกันภายใต้ Clean Architecture ได้อย่างไร โดยใช้ business logic ชุดเดียวกันทั้งหมด

---

## 🧱 Tech Stack

| Layer        | Technology                                  |
| ------------ | ------------------------------------------- |
| Backend      | Go 1.22+ · Fiber · GORM · PostgreSQL        |
| API          | REST · GraphQL                              |
| Realtime     | WebSocket                                   |
| Frontend     | React 19 · TypeScript · Vite · Tailwind CSS |
| Architecture | Clean Architecture                          |

---

## โครงสร้างโปรเจค
```text
                    React Frontend
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
      REST             GraphQL          WebSocket
        │                  │                  │
        └──────────────────┴──────────────────┘
                           │
                        Usecase
                    (Business Logic)
                           │
                      Repository
                           │
                      PostgreSQL
```

Design หลักของโปรเจคคือ

> REST, GraphQL และ WebSocket เป็นเพียง delivery layer คนละรูปแบบ แต่ทั้งหมดเรียก usecase ชุดเดียวกัน

ทำให้ business logic ไม่กระจายอยู่หลายที่ และสามารถเปลี่ยน transport layer ได้โดยไม่กระทบ core ของระบบ

---

## Why REST?

REST ถูกใช้กับ operation ที่เป็นลักษณะ command เช่น

* Create User
* Create Room
* Send Message

เนื่องจาก request และ response มีรูปแบบค่อนข้างตายตัวอยู่แล้ว การใช้ REST ทำให้ implementation ตรงไปตรงมาและ debug ได้ง่าย

---

## Why GraphQL?

ในหน้าแชท เราไม่ได้ต้องการแค่ข้อความ แต่ต้องการข้อมูลของผู้ส่งด้วย

ตัวอย่างเช่น หากดึง message history ผ่าน REST

```text
GET /rooms/1/messages
```

จะได้ข้อมูลประมาณ

```json
{
  "user_id": 1,
  "content": "hello"
}
```

จากนั้น frontend ต้องนำ user_id ไปดึงข้อมูล user เพิ่มอีกครั้งเพื่อแสดง username

```text
fetch messages    // เรียก API เพื่อดึงข้อความล่าสุดในห้อง (id, room_id, user_id, content, created_at)
      ↓
extract user ids // ลูปหาข้อมูลทั้งหมดเพื่อรวบรวม user_id ที่ไม่ซ้ำกัน
      ↓
fetch users      // เรียก API อีกรอบเพื่อดึงข้อมูลผู้ใช้ตาม id ที่รวบรวมได้
      ↓
create map      // สร้าง map เช่น {1: "alice", 2: "bob"} เพื่อให้ lookup ได้ง่าย
      ↓
merge data      // นำข้อมูล user มา merge กับ message แต่ละรายการ
      ↓
render แสดงผล
```

GraphQL ช่วยลดขั้นตอนเหล่านี้ โดย client สามารถระบุ field ที่ต้องการได้ตั้งแต่แรก

```graphql
query {
  messages(roomId: 1, limit: 25) {
    content
    createdAt

    user {
      username
    }
  }
}
```

ผลลัพธ์ที่ได้คือ

```json
[
  {
    "content": "สวัสดี",
    "createdAt": "2026-06-18T10:00:00Z",
    "user": {
      "username": "alice"
    }
  }
]
```

frontend สามารถนำข้อมูลไป render ได้ทันที

```text
GraphQL

fetch
  ↓
map
  ↓
done
```
ไม่ต้องแก้หรือสร้าง endpoint เยอะ 
Frontend ไม่ต้อง fetch หลายรอบ
Backend ไม่ต้อง implement การ join ให้ยุ่งยาก หากเปลี่ยน schema


---

## Why WebSocket?

REST และ GraphQL เป็น request-response protocol

client ต้องเป็นฝ่ายถามก่อนทุกครั้ง

แต่สำหรับระบบ chat เราไม่อยากให้ browser คอย polling ตลอดเวลาเพื่อเช็คว่ามีข้อความใหม่หรือไม่

หากใช้ polling ฝั่ง browser จะต้องคอยยิง request ซ้ำ ๆ เพื่อถามว่ามีข้อความใหม่หรือไม่ แม้ในช่วงที่ไม่มีใครส่งข้อความเลยก็ตาม ทำให้เกิด request จำนวนมากโดยไม่จำเป็น 
WebSocket ถูกเพิ่มเข้ามาเพื่อให้ server สามารถ push event ไปหา client ได้โดยตรง

โครงสร้าง

```text
ส่งข้อความ
    │
    ▼

REST Endpoint

    ▼

Usecase

    ▼

PostgreSQL

    ▼

MessageNotifier

    ▼

WebSocket Hub

    ▼

Broadcast
```

เมื่อมีการส่งข้อความใหม่

1. REST handler รับ request
2. Usecase บันทึกข้อมูลลง database
3. Usecase เรียก MessageNotifier
4. WebSocket Hub broadcast ไปยังทุก client ที่อยู่ในห้องเดียวกัน

ดังนั้น usecase ไม่จำเป็นต้องรู้ว่าปลายทางคือ WebSocket หรือเทคโนโลยีอื่น

---

## Dependency Inversion

Usecase ไม่ได้อ้างอิง WebSocket โดยตรง

แต่ทำงานผ่าน interface

```go
type MessageNotifier interface {
    BroadcastMessage(message Message)
}
```

ทำให้ business logic สามารถถูก test ได้ง่ายขึ้น และไม่ผูกติดกับ implementation รายละเอียดของ transport layer

---

## 🔌 REST API

| Method | Path                      | Description          |
| ------ | ------------------------- | -------------------- |
| POST   | `/api/users`              | Create or fetch user |
| GET    | `/api/rooms`              | List rooms           |
| POST   | `/api/rooms`              | Create room          |
| GET    | `/api/rooms/:id/messages` | Message history      |
| POST   | `/api/rooms/:id/messages` | Send message         |

