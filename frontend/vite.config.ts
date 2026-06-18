import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // โหลดค่าจากไฟล์ .env (prefix '' = อ่านทุกตัว ไม่จำกัดแค่ VITE_)
  // BACKEND_URL ใช้แค่ฝั่ง dev proxy (รันบน Node) ไม่ได้ส่งไป browser จึงไม่ต้องมี prefix VITE_
  const env = loadEnv(mode, process.cwd(), '')
  const backendUrl = env.BACKEND_URL || 'http://localhost:8085'

  return {
    plugins: [react(), tailwindcss()],
    server: {
      // อนุญาตให้เข้าผ่าน host ของ ngrok ได้ (.ngrok-free.app = ทุก subdomain ของ ngrok)
      // ถ้าใช้ tunnel เจ้าอื่นก็เพิ่ม host เข้าไปใน list นี้ หรือใช้ true เพื่ออนุญาตทุก host
      allowedHosts: ['.ngrok-free.app'],
      proxy: {
        // ทุก request ที่ขึ้นต้นด้วย /api จะถูก proxy ไป backend (เปลี่ยน port ได้ที่ BACKEND_URL ใน .env)
        // ทำให้ frontend เรียก /api/rooms ได้เลยโดยไม่ต้องเขียน full URL + ไม่ติด CORS ตอน dev
        '/api': {
          target: backendUrl,
          changeOrigin: true,
        },
        // GraphQL endpoint — proxy ให้เป็น same-origin เหมือน /api (เลี่ยง CORS ตอน dev)
        '/graphql': {
          target: backendUrl,
          changeOrigin: true,
        },
        // WebSocket endpoint — ws: true ให้ proxy ส่งต่อการ upgrade เป็น WebSocket ได้
        '/ws': {
          target: backendUrl,
          changeOrigin: true,
          ws: true,
        },
      },
    },
  }
})
