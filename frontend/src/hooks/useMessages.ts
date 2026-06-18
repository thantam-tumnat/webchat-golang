import { useEffect } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { listMessages, sendMessage } from '../api/endpoints'
import type { Message, Paginated } from '../types'

// useMessages: โหลดประวัติครั้งแรกด้วย GraphQL แล้ว "ฟัง" ข้อความใหม่ผ่าน WebSocket
// (มาแทน polling เดิมที่ยิงถามทุก 3 วิ) — server จะ push ข้อความใหม่มาหาเองตอนมีจริง
export function useMessages(roomId: number) {
  const qc = useQueryClient()

  const query = useQuery({
    queryKey: ['messages', roomId],
    queryFn: () => listMessages(roomId, 1, 50),
    enabled: roomId > 0,
    // กันเหนียว: เผื่อ WebSocket หลุด (เน็ตสะดุด) ยัง refetch เป็นพัก ๆ ให้ข้อความไม่ตกหล่น (fallback)
    refetchInterval: 15000,
  })

  // เปิด WebSocket ฟังข้อความสดของห้องนี้ (เปิดใหม่ทุกครั้งที่เปลี่ยนห้อง)
  useEffect(() => {
    if (roomId <= 0) return

    // เลือก protocol ตามหน้าเว็บ: http -> ws, https -> wss
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const ws = new WebSocket(
      `${proto}://${window.location.host}/ws/rooms/${roomId}`,
    )

    // มีข้อความใหม่เข้ามา → ยัดเข้า cache ของ React Query ตรง ๆ (ไม่ต้อง refetch)
    ws.onmessage = (event) => {
      const incoming: Message = JSON.parse(event.data)
      qc.setQueryData<Paginated<Message>>(['messages', roomId], (old) => {
        if (!old) return old
        // กันซ้ำ: ถ้ามี id นี้อยู่แล้ว (เช่นเพิ่งมาจาก refetch) ไม่ต้องเพิ่มอีก
        if (old.data.some((m) => m.id === incoming.id)) return old
        // cache เรียงใหม่->เก่า (DESC) → ข้อความใหม่ไปไว้หน้าสุด (ChatPage จะ reverse แสดงล่างสุดเอง)
        return { ...old, data: [incoming, ...old.data], total: old.total + 1 }
      })
    }

    // ปิด connection เมื่อออกจากห้อง/หน้าเปลี่ยน กัน connection ค้าง
    return () => ws.close()
  }, [roomId, qc])

  return query
}

// useSendMessage ส่งข้อความผ่าน REST เหมือนเดิม
// ข้อความที่ส่งสำเร็จจะถูก backend broadcast กลับมาทาง WebSocket ให้ทุกคน (รวมตัวเอง)
// คง invalidate ไว้เป็น fallback เผื่อ WebSocket ไม่ติด — dedupe ฝั่ง onmessage กันซ้ำให้แล้ว
export function useSendMessage(roomId: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ userId, content }: { userId: number; content: string }) =>
      sendMessage(roomId, userId, content),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['messages', roomId] })
    },
  })
}
