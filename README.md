# 💬 Real-time Chat App (Go + React)

> 📖 ดูโครงสร้างโปรเจค, request flow และขั้นตอนการติดตั้งได้ที่ [workflow_explanation.md](workflow_explanation.md)

---

## 🧱 Tech Stack

| ชั้น (Layer) | เทคโนโลยี |
|-------|------------|
| Backend | Go 1.22+ · Fiber v2 · GORM · PostgreSQL |
| API | REST (`/api`) · GraphQL (`/graphql`, graphql-go) |
| Frontend | React 19 · TypeScript · Vite · Tailwind CSS v4 |
| สถาปัตยกรรม | Clean Architecture (domain / usecase / repository / delivery) |

---

## 🔌 API Endpoints

### REST (Phase 1)

| Method | Path | คำอธิบาย |
|--------|------|-------------|
| GET | `/health` | ตรวจสถานะระบบ (health check) |
| POST | `/api/users` | สร้างหรือดึง user — `{ "username": "..." }` |
| GET | `/api/rooms` | ดึงรายชื่อห้องทั้งหมด |
| POST | `/api/rooms` | สร้างห้องใหม่ — `{ "name": "..." }` |
| GET | `/api/rooms/:id/messages?page=1&limit=20` | ดึงข้อความในห้องแบบแบ่งหน้า (pagination) |
| POST | `/api/rooms/:id/messages` | ส่งข้อความ — `{ "user_id": 1, "content": "..." }` |

### GraphQL (`POST /graphql`)

REST คืนข้อมูล **รูปร่างตายตัวต่อ endpoint** ส่วน GraphQL เพิ่มเข้ามาเพื่อให้ client
ขอ **เฉพาะ field ที่ต้องการในคำขอเดียว** — ไม่ over-fetch และไม่ต้องยิงหลายรอบ

**Design หลัก — GraphQL เป็นแค่ delivery layer อีกตัว ที่เรียก usecase ชุดเดียวกับ REST**
business logic จึงไม่ถูกเขียนซ้ำ ซึ่งคือหัวใจของ Clean Architecture:

```
REST handler   ─┐
                ├─► usecase (business logic) ─► repository ─► PostgreSQL
GraphQL resolver ┘
```

> ตัวอย่าง: resolver ของ `rooms` คือ `return r.roomUC.List(ctx)` ตรง ๆ — เป็น usecase
> ตัวเดียวกับที่ REST handler ใช้ การเพิ่ม GraphQL จึงไม่ต้องแตะชั้น domain หรือ repository เลย

**Schema**

| ประเภท | Field | Args | คำอธิบาย |
|------|-------|------|-------------|
| Query | `rooms` | — | ดึงรายชื่อห้องทั้งหมด |
| Query | `messages` | `roomId!`, `page`, `limit` | ดึงข้อความในห้องแบบแบ่งหน้า |
| Mutation | `createUser` | `username!` | สร้างหรือดึง user |
| Mutation | `createRoom` | `name!` | สร้างห้อง |
| Mutation | `sendMessage` | `roomId!`, `userId!`, `content!` | ส่งข้อความ |

**ตัวอย่าง** — ดึงรายชื่อห้องและข้อความล่าสุดในห้อง โดยเลือกเฉพาะ field ที่ UI ต้องใช้:

```graphql
query {
  rooms { id name }
  messages(roomId: 1, limit: 5) {
    content
    user { username }
  }
}
```

> เปิดใช้ **GraphiQL playground** ไว้ — เข้า `GET /graphql` ในเบราว์เซอร์เพื่อดู schema
> และลองรัน query แบบ interactive ได้ทันที

#### ตัวอย่างการใช้งาน

**1) Query — ดึงข้อความในห้องพร้อมชื่อคนส่ง (ใช้ตัวแปร)**

```graphql
query Messages($roomId: Int!, $limit: Int) {
  messages(roomId: $roomId, limit: $limit) {
    content
    user { username }
  }
}
```
```json
// variables
{ "roomId": 1, "limit": 2 }

// response
{
  "data": {
    "messages": [
      { "content": "สวัสดีทุกคน", "user": { "username": "alice" } },
      { "content": "หวัดดีครับ",   "user": { "username": "bob" } }
    ]
  }
}
```

> จุดสำคัญ: ได้ทั้งข้อความและ `user.username` **ในคำขอเดียว** — ถ้าเป็น REST จะได้แค่ `user_id`
> แล้วต้องยิงขอ user เพิ่มอีกรอบ

**2) Mutation — ส่งข้อความ แล้วเลือกรับเฉพาะ field ที่ต้องใช้กลับมา**

```graphql
mutation {
  sendMessage(roomId: 1, userId: 1, content: "สวัสดีทุกคน") {
    id
    content
    createdAt
  }
}
```
```json
// response
{
  "data": {
    "sendMessage": { "id": 42, "content": "สวัสดีทุกคน", "createdAt": "2026-06-18T10:30:00Z" }
  }
}
```

> เรียกผ่าน HTTP ก็ได้: `POST /graphql` ด้วย body `{ "query": "...", "variables": { ... } }`

### REST vs GraphQL — ใช้ตัวไหนตอนไหน

โปรเจคนี้ **ไม่ได้ใช้ GraphQL ทุกที่** แต่เลือกใช้ตามความเหมาะสมของแต่ละงาน:

| งาน | ใช้ | เหตุผล |
|------|------|--------|
| โหลดข้อความในห้อง | **GraphQL** | ต้องการ message พร้อมข้อมูล `user` (ชื่อคนส่ง) ที่สัมพันธ์กันในคำขอเดียว — เลี่ยงการยิงขอ user ซ้ำ |
| สร้าง user / ห้อง, ดึงรายชื่อห้อง, ส่งข้อความ | REST | เป็น CRUD ตรงไปตรงมา ไม่มี relation ซับซ้อน |
