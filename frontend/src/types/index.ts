// types เหล่านี้ต้องตรงกับ JSON ที่ backend ส่งกลับมา (ดู struct ใน backend/internal/domain)

export interface User {
  id: number
  username: string
  created_at: string
}

export interface Room {
  id: number
  name: string
  created_at: string
}

export interface Message {
  id: number
  room_id: number
  user_id: number
  content: string
  created_at: string
  user?: User // backend preload มาให้ เพื่อแสดงชื่อผู้ส่ง
}

// รูปแบบ response แบบแบ่งหน้า (ตรงกับ PaginatedResponse ฝั่ง backend)
export interface Paginated<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

// รูปแบบ error จาก backend
export interface ApiError {
  code: string
  message: string
  details?: unknown
}
