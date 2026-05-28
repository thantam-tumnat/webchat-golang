import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // อนุญาตให้เข้าผ่าน host ของ ngrok ได้ (.ngrok-free.app = ทุก subdomain ของ ngrok)
    // ถ้าใช้ tunnel เจ้าอื่นก็เพิ่ม host เข้าไปใน list นี้ หรือใช้ true เพื่ออนุญาตทุก host
    allowedHosts: ['.ngrok-free.app'],
    proxy: {
      // ทุก request ที่ขึ้นต้นด้วย /api จะถูกส่งต่อไป backend ที่ port 8080
      // ทำให้ frontend เรียก /api/rooms ได้เลยโดยไม่ต้องเขียน full URL + ไม่ติด CORS ตอน dev
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
