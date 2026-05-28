import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { listMessages, sendMessage } from '../api/endpoints'

// useMessages ดึงข้อความในห้อง + polling ทุก 3 วินาที (Phase 1 ยังไม่มี WebSocket)
// refetchInterval = 3000 -> React Query จะ refetch อัตโนมัติทุก 3 วิ ทำให้เห็นข้อความใหม่จากคนอื่น
export function useMessages(roomId: number) {
  return useQuery({
    queryKey: ['messages', roomId],
    queryFn: () => listMessages(roomId, 1, 50),
    refetchInterval: 3000,
    enabled: roomId > 0,
  })
}

// useSendMessage ส่งข้อความ แล้ว refetch ข้อความในห้องนั้นทันที
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
