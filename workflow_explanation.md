# 🧭 Workflow Explanation

A full-stack chat application built in **phases** — starting with a working CRUD foundation,
then gradually layering on features (Auth, WebSocket, Redis, Deploy) phase by phase.

> **Current status: Phase 1 — CRUD + Polling** ✅

---

## 📂 Project Structure

### Backend (`backend/`)

ตารางเทียบโครงสร้าง RESTAPI กับ GraphQL
![alt text](image.png)
```
main.go                Entry point — loads config, connects to DB, wires up
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

  delivery/graphql/    ★ Another delivery layer — GraphQL, reuses the same usecases.
    types.go           GraphQL object types (User / Room / Message) mapped to domain
    schema.go          Query + Mutation definitions
    resolver.go        Resolvers — call the existing usecases (no logic duplicated)
    handler.go         GraphQL handler + GraphiQL playground, wrapped for Fiber

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

GraphQL takes the **same path** from the usecase down — only the entry point differs:
```
POST /graphql  mutation { sendMessage(...) }
  → graphql handler
  → resolver.sendMessage()
  → message_usecase.Send()      ← same business logic as REST
  → message_repo.Create()
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
go run .
```
- On success you'll see `🚀 server running at http://localhost:8080`
- Tables are created automatically (AutoMigrate)
- Quick check: open http://localhost:8080/health — you should see `{"status":"ok"}`
- GraphQL playground (GraphiQL): open http://localhost:8080/graphql in the browser

### 3) Run the Frontend
```bash
cd frontend
npm install
npm run dev
```
- Open http://localhost:5173
