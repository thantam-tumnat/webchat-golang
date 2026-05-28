import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createRoom, listRooms } from '../api/endpoints'

// useRooms ดึงรายการห้องทั้งหมด (React Query จัดการ cache + loading + error ให้)
export function useRooms() {
  return useQuery({
    queryKey: ['rooms'],
    queryFn: listRooms,
  })
}

// useCreateRoom สร้างห้องใหม่ แล้ว invalidate cache เพื่อให้ list โหลดใหม่อัตโนมัติ
export function useCreateRoom() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (name: string) => createRoom(name),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['rooms'] })
    },
  })
}
