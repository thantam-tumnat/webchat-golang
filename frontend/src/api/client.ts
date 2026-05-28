import axios from 'axios'

// axios instance กลางของแอป
// baseURL = /api -> ตอน dev จะถูก vite proxy ส่งต่อไป backend (ดู vite.config.ts)
export const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
})

// หมายเหตุ: Phase 2 (Auth) จะเพิ่ม interceptor ตรงนี้
// เพื่อแนบ JWT token ใน header และ refresh token อัตโนมัติ
