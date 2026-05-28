import { api } from './client'
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
export async function listMessages(
  roomId: number,
  page = 1,
  limit = 20,
): Promise<Paginated<Message>> {
  const { data } = await api.get<Paginated<Message>>(
    `/rooms/${roomId}/messages`,
    { params: { page, limit } },
  )
  return data
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
