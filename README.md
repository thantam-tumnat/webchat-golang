# 💬 Real-time Chat App (Go + React)

A full-stack chat application built in **phases** — starting with a working CRUD foundation,
then gradually layering on features (Auth, WebSocket, Redis, Deploy) phase by phase.

> **Current status: Phase 1 — CRUD + Polling** ✅

---

## 🧱 Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+ · Fiber v2 · GORM · PostgreSQL |
| Frontend | React 19 · TypeScript · Vite · Tailwind CSS v4 |
| State / Data | Zustand · TanStack Query (React Query) · Axios |
| Forms | React Hook Form · Zod |
| Architecture | Clean Architecture (domain / usecase / repository / delivery) |

---

## 📂 Project Structure

### Backend (`backend/`)

```
cmd/api/main.go        Entry point — loads config, connects to DB, wires up
                       every layer (dependency injection), starts the server,
                       and handles graceful shutdown.

internal/
  config/config.go     Loads configuration from .env / environment variables.

  domain/              ★ The core — no dependency on any framework or DB.
    user.go            User entity + UserRepository interface
    room.go            Room entity + RoomRepository interface
    message.go         Message entity + MessageRepository interface
    errors.go          AppError — the system's standard error shape {code, message}

  usecase/             ★ Business logic — only knows the repository interfaces.
    user_usecase.go    Create / fetch user (idempotent by username)
    room_usecase.go    List / create rooms
    message_usecase.go Send / fetch messages (verifies room + user exist)

  repository/          ★ Real implementations that talk to the DB via GORM.
    user_repo.go       (the only place that "knows about" GORM)
    room_repo.go
    message_repo.go

  delivery/http/       ★ The outside-world layer (HTTP) — Fiber lives here.
    router.go          Registers middleware (CORS / logger / recover) + routes
    response.go        Central error handler + validator + response format
    user_handler.go    Parses requests, validates input, calls the usecase
    room_handler.go
    message_handler.go

  infrastructure/
    database/postgres.go  Opens the GORM connection + runs AutoMigrate

migrations/            SQL migration files (for golang-migrate in a later phase)
docker-compose.yml     Runs PostgreSQL + Adminer
.env.example           Example config — copy this to .env
Makefile               Shortcut commands (db-up / run / build)
```

**Request flow** (example: sending a message):
```
POST /api/rooms/1/messages
  → router.go (passes through middleware)
  → message_handler.Send()      parse + validate
  → message_usecase.Send()      business logic
  → message_repo.Create()       INSERT into DB
  → return JSON response
```

### Frontend (`frontend/src/`)

```
main.tsx               Wraps the app in the React Query provider.
App.tsx                Router — UsernameGate wraps the routes / and /rooms/:id

types/index.ts         TypeScript types mirroring the backend's JSON shapes
api/client.ts          Axios instance (baseURL = /api)
api/endpoints.ts       One function per API endpoint

stores/userStore.ts    Stores the current user via Zustand + localStorage

hooks/useRooms.ts      React Query: list / create rooms
hooks/useMessages.ts   React Query: fetch messages (polls every 3s) + send

components/
  UsernameGate.tsx     Forces the user to set a username before entering
pages/
  RoomsPage.tsx        Room list + create-room form
  ChatPage.tsx         Chat view — renders messages and sends new ones
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+
- Docker Desktop (for running PostgreSQL)

### 1) Start PostgreSQL with Docker
```bash
cd backend
docker compose up -d
```
- PostgreSQL runs on `localhost:5432`
- Adminer (web-based DB viewer) is at http://localhost:8081
  (System: PostgreSQL · Server: `postgres` · User/Pass/DB: `chatapp`)

### 2) Run the Backend
```bash
cd backend
cp .env.example .env        # Windows: copy .env.example .env
go run ./cmd/api
```
- On success you'll see `🚀 server running at http://localhost:8080`
- Tables are created automatically (AutoMigrate)
- Quick check: open http://localhost:8080/health — you should see `{"status":"ok"}`

### 3) Run the Frontend
```bash
cd frontend
npm install
npm run dev
```
- Open http://localhost:5173


## 🔌 API Endpoints (Phase 1)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/api/users` | Create or fetch a user — `{ "username": "..." }` |
| GET | `/api/rooms` | List all rooms |
| POST | `/api/rooms` | Create a room — `{ "name": "..." }` |
| GET | `/api/rooms/:id/messages?page=1&limit=20` | Paginated messages in a room |
| POST | `/api/rooms/:id/messages` | Send a message — `{ "user_id": 1, "content": "..." }` |
