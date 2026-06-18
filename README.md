# 💬 Real-time Chat App (Go + React)

> 📖 For project structure, request flow, and setup steps, see [workflow_explanation.md](workflow_explanation.md).

---

## 🧱 Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+ · Fiber v2 · GORM · PostgreSQL |
| API | REST (`/api`) · GraphQL (`/graphql`, graphql-go) |
| Frontend | React 19 · TypeScript · Vite · Tailwind CSS v4 |
| State / Data | Zustand · TanStack Query (React Query) · Axios |
| Forms | React Hook Form · Zod |
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

### GraphQL (`/graphql`)

Open `http://localhost:8080/graphql` for the GraphiQL playground.

```graphql
type Query {
  rooms: [Room]
  messages(roomId: Int!, page: Int, limit: Int): [Message]
}

type Mutation {
  createUser(username: String!): User
  createRoom(name: String!): Room
  sendMessage(roomId: Int!, userId: Int!, content: String!): Message
}
```

Example — fetch messages with the sender's username in a single request:
```graphql
query {
  messages(roomId: 1) {
    content
    user { username }
  }
}
```
