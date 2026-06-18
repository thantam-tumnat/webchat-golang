import axios from 'axios'
import { useUserStore } from '../stores/userStore'

// axios instance กลางของแอป
// baseURL = /api -> ตอน dev จะถูก vite proxy ส่งต่อไป backend (ดู vite.config.ts)
export const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
})

// Response interceptor: ดักทุก response ก่อนถึงโค้ดส่วนอื่น
// ถ้า backend ตอบ USER_NOT_FOUND แปลว่า user ที่จำไว้ใน localStorage
// ไม่มีใน DB แล้ว (เช่นหลังล้าง DB) -> เคลียร์ทิ้งให้อัตโนมัติ
// store เปลี่ยนเป็น user=null -> UsernameGate เด้งขึ้นให้ตั้งชื่อใหม่เอง
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.data?.code === 'USER_NOT_FOUND') {
      useUserStore.getState().logout()
    }
    return Promise.reject(err)
  },
)

// หมายเหตุ: Phase 2 (Auth) จะเพิ่ม interceptor ตรงนี้
// เพื่อแนบ JWT token ใน header และ refresh token อัตโนมัติ
