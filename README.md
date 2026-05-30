# 💬 Real-time Chat App (Go + React)

โปรเจกต์แชทแบบ full-stack ทำเป็น **phase** — เริ่มจาก CRUD พื้นฐานที่รันได้จริง
แล้วค่อยๆ เพิ่ม feature (Auth, WebSocket, Redis, Deploy) ในแต่ละ phase

> **สถานะปัจจุบัน: Phase 1 — CRUD + Polling** ✅

---

## 🧱 Tech Stack

| ส่วน | เทคโนโลยี |
|------|-----------|
| Backend | Go 1.22+ · Fiber v2 · GORM · PostgreSQL |
| Frontend | React 19 · TypeScript · Vite · Tailwind CSS v4 |
| State / Data | Zustand · TanStack Query (React Query) · Axios |
| Form | React Hook Form · Zod |
| Architecture | Clean Architecture (domain / usecase / repository / delivery) |

---

## 📂 โครงสร้างโปรเจกต์ + แต่ละไฟล์ทำอะไร

### Backend (`backend/`)

```
cmd/api/main.go        จุดเริ่มของโปรแกรม — โหลด config, ต่อ DB, ประกอบ (wiring)
                       ทุก layer เข้าด้วยกัน (Dependency Injection), รัน server,
                       graceful shutdown

internal/
  config/config.go     โหลดค่า config จาก .env / environment variables

  domain/              ★ แกนกลาง — ไม่ขึ้นกับ framework/DB ใดๆ
    user.go            entity User + interface UserRepository
    room.go            entity Room + interface RoomRepository
    message.go         entity Message + interface MessageRepository
    errors.go          AppError (error มาตรฐานของระบบ {code, message})

  usecase/             ★ business logic — รู้จักแค่ interface ของ repository
    user_usecase.go    สร้าง/หา user (idempotent ตาม username)
    room_usecase.go    list / สร้างห้อง
    message_usecase.go ส่ง/ดึงข้อความ (ตรวจว่าห้อง+user มีจริง)

  repository/          ★ implementation จริงที่คุย DB ด้วย GORM
    user_repo.go       (เป็นที่เดียวที่ "รู้จัก" GORM)
    room_repo.go
    message_repo.go

  delivery/http/       ★ ชั้นติดต่อโลกภายนอก (HTTP) — Fiber อยู่ตรงนี้
    router.go          ลงทะเบียน middleware (CORS/logger/recover) + routes
    response.go        custom error handler กลาง + validator + รูปแบบ response
    user_handler.go    รับ request, parse, validate, เรียก usecase
    room_handler.go
    message_handler.go

  infrastructure/
    database/postgres.go  เปิด connection GORM + AutoMigrate

migrations/            SQL migration files (ไว้ใช้กับ golang-migrate ใน phase หลัง)
docker-compose.yml     รัน PostgreSQL + Adminer
.env.example           ตัวอย่างค่า config (คัดลอกเป็น .env)
Makefile               คำสั่งลัด (db-up / run / build)
```

**ลำดับการไหลของ request** (เช่น ส่งข้อความ):
```
POST /api/rooms/1/messages
  → router.go (ผ่าน middleware)
  → message_handler.Send()      parse + validate
  → message_usecase.Send()      business logic
  → message_repo.Create()       INSERT ลง DB
  → ส่ง JSON กลับ
```

### Frontend (`frontend/src/`)

```
main.tsx               ใส่ Provider ของ React Query ครอบทั้งแอป
App.tsx                router — UsernameGate ครอบ + เส้นทาง / และ /rooms/:id

types/index.ts         TypeScript types ที่ตรงกับ JSON ของ backend
api/client.ts          axios instance (baseURL = /api)
api/endpoints.ts        ฟังก์ชันเรียก API แต่ละตัว

stores/userStore.ts    เก็บ user ปัจจุบันด้วย Zustand + localStorage

hooks/useRooms.ts      React Query: list/สร้างห้อง
hooks/useMessages.ts   React Query: ดึงข้อความ (polling ทุก 3 วิ) + ส่งข้อความ

components/
  UsernameGate.tsx     บังคับตั้ง username ก่อนเข้าใช้งาน
pages/
  RoomsPage.tsx        หน้ารายการห้อง + สร้างห้อง
  ChatPage.tsx         หน้าแชท — แสดงข้อความ + ส่งข้อความ
```

---

## 🚀 วิธีรัน (ทำทีละขั้น)

### สิ่งที่ต้องมี
- Go 1.22+
- Node.js 18+
- Docker Desktop (สำหรับรัน PostgreSQL)

### 1) เปิด PostgreSQL ด้วย Docker
```bash
cd backend
docker compose up -d
```
- PostgreSQL จะรันที่ `localhost:5432`
- Adminer (หน้าเว็บดู DB) ที่ http://localhost:8081
  (System: PostgreSQL, Server: `postgres`, User/Pass/DB: `chatapp`)

### 2) รัน Backend
```bash
cd backend
cp .env.example .env        # windows: copy .env.example .env
go run ./cmd/api
```
- ถ้าสำเร็จจะเห็น `🚀 server กำลังรันที่ http://localhost:8080`
- ตารางถูกสร้างให้อัตโนมัติ (AutoMigrate)
- ทดสอบ: เปิด http://localhost:8080/health ควรเห็น `{"status":"ok"}`

### 3) รัน Frontend
```bash
cd frontend
npm install
npm run dev
```
- เปิด http://localhost:5173


## 🔌 API Endpoints (Phase 1)

| Method | Path | คำอธิบาย |
|--------|------|----------|
| GET | `/health` | health check |
| POST | `/api/users` | สร้าง/หา user `{ "username": "..." }` |
| GET | `/api/rooms` | รายการห้องทั้งหมด |
| POST | `/api/rooms` | สร้างห้อง `{ "name": "..." }` |
| GET | `/api/rooms/:id/messages?page=1&limit=20` | ข้อความในห้อง (แบ่งหน้า) |
| POST | `/api/rooms/:id/messages` | ส่งข้อความ `{ "user_id": 1, "content": "..." }` |

