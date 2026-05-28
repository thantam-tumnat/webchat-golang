import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../types'

// userStore เก็บข้อมูล user ปัจจุบัน (Phase 1 ยังไม่มี login จริง แค่ตั้ง username)
// ใช้ persist middleware เพื่อเก็บลง localStorage -> refresh แล้วยังจำได้
interface UserState {
  user: User | null
  setUser: (user: User) => void
  logout: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      user: null,
      setUser: (user) => set({ user }),
      logout: () => set({ user: null }),
    }),
    { name: 'chatapp-user' }, // key ใน localStorage
  ),
)
