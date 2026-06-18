# 💬 Real-time Chat App (Go + React)

> 📖 For project structure, request flow, and setup steps, see [workflow_explanation.md](workflow_explanation.md).

---

## 🧱 Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+ · Fiber v2 · GORM · PostgreSQL |
| API | REST (`/api`) · GraphQL (`/graphql`, graphql-go) |
| Frontend | React 19 · TypeScript · Vite · Tailwind CSS v4 |
| Architecture | Clean Architecture (domain / usecase / repository / delivery) |

---

## 🔌 API Endpoints

### REST (Phase 1)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/api/users` | Create or fetch a user — `{ "username": "..." }` |
| GET | `/api/rooms` | List all rooms |
| POST | `/api/rooms` | Create a room — `{ "name": "..." }` |
| GET | `/api/rooms/:id/messages?page=1&limit=20` | Paginated messages in a room |
| POST | `/api/rooms/:id/messages` | Send a message — `{ "user_id": 1, "content": "..." }` |

### GraphQL (`POST /graphql`)

REST returns a **fixed shape per endpoint**. GraphQL is added on top so the client can
ask for **exactly the fields it needs in a single request** — no over-fetching, no extra round trips.

**Key design — GraphQL is just another delivery layer; it calls the same usecases as REST.**
Business logic is never duplicated, which is the whole point of Clean Architecture:

```
REST handler   ─┐
                ├─► usecase (business logic) ─► repository ─► PostgreSQL
GraphQL resolver ┘
```

> e.g. the `rooms` resolver is literally `return r.roomUC.List(ctx)` — the same usecase the REST
> handler uses. Adding GraphQL didn't touch the domain or repository layer at all.

**Schema**

| Kind | Field | Args | Description |
|------|-------|------|-------------|
| Query | `rooms` | — | List all rooms |
| Query | `messages` | `roomId!`, `page`, `limit` | Paginated messages in a room |
| Mutation | `createUser` | `username!` | Create or fetch a user |
| Mutation | `createRoom` | `name!` | Create a room |
| Mutation | `sendMessage` | `roomId!`, `userId!`, `content!` | Send a message |

**Example** — fetch rooms and a room's latest messages, picking only the fields the UI needs:

```graphql
query {
  rooms { id name }
  messages(roomId: 1, limit: 5) {
    content
    user { username }
  }
}
```

> A **GraphiQL playground** is enabled — open `GET /graphql` in the browser to explore the
> schema and run queries interactively.

