import { api } from './client'
import { gql } from './graphql'
import type { Message, Paginated, Room, User } from '../types'

// --- Users ---
export async function createUser(username: string): Promise<User> {
  const { data } = await api.post<User>('/users', { username })
  return data
}

// --- Rooms ---
export async function listRooms(): Promise<Room[]> {
  const { data } = await api.get<Room[]>('/rooms')
  return data
}

export async function createRoom(name: string): Promise<Room> {
  const { data } = await api.post<Room>('/rooms', { name })
  return data
}

// --- Messages ---
// ★ ดึงข้อความผ่าน GraphQL (ฝั่งอ่านที่มี relation: message + user ในคำขอเดียว)
//   ใช้ alias (เช่น user_id: userId) เปลี่ยนชื่อ field ของ GraphQL (camelCase)
//   ให้กลับมาเป็น snake_case ตรงกับ type Message เดิม -> component ไม่ต้องแก้
const MESSAGES_QUERY = `
  query Messages($roomId: Int!, $page: Int, $limit: Int) {
    messages(roomId: $roomId, page: $page, limit: $limit) {
      id
      room_id: roomId
      user_id: userId
      content
      created_at: createdAt
      user {
        id
        username
      }
    }
  }
`

export async function listMessages(
  roomId: number,
  page = 1,
  limit = 20,
): Promise<Paginated<Message>> {
  const data = await gql<{ messages: Message[] }>(MESSAGES_QUERY, {
    roomId,
    page,
    limit,
  })
  // ห่อกลับเป็นรูปแบบ Paginated เดิม เพื่อให้ useMessages / ChatPage ใช้ต่อได้ไม่ต้องแก้
  return { data: data.messages, total: data.messages.length, page, limit }
}

export async function sendMessage(
  roomId: number,
  userId: number,
  content: string,
): Promise<Message> {
  const { data } = await api.post<Message>(`/rooms/${roomId}/messages`, {
    user_id: userId,
    content,
  })
  return data
}
