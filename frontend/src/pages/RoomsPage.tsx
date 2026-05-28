import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useCreateRoom, useRooms } from '../hooks/useRooms'
import { useUserStore } from '../stores/userStore'

const schema = z.object({
  name: z.string().min(1, 'ใส่ชื่อห้อง').max(100, 'ชื่อยาวเกินไป'),
})
type FormValues = z.infer<typeof schema>

export function RoomsPage() {
  const { user, logout } = useUserStore()
  const { data: rooms, isLoading, isError } = useRooms()
  const createRoom = useCreateRoom()

  const { register, handleSubmit, reset, formState: { errors } } =
    useForm<FormValues>({ resolver: zodResolver(schema) })

  return (
    <div className="mx-auto min-h-screen max-w-2xl p-4">
      <header className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold text-slate-800">💬 ห้องแชท</h1>
        <div className="flex items-center gap-3 text-sm">
          <span className="text-slate-500">
            สวัสดี <b className="text-slate-700">{user?.username}</b>
          </span>
          <button
            onClick={logout}
            className="rounded-md border border-slate-300 px-3 py-1 text-slate-600 hover:bg-slate-100"
          >
            ออก
          </button>
        </div>
      </header>

      {/* ฟอร์มสร้างห้อง */}
      <form
        onSubmit={handleSubmit((v) => {
          createRoom.mutate(v.name, { onSuccess: () => reset() })
        })}
        className="mb-6 flex gap-2"
      >
        <input
          {...register('name')}
          placeholder="ชื่อห้องใหม่..."
          className="flex-1 rounded-lg border border-slate-300 px-4 py-2 outline-none focus:border-blue-500"
        />
        <button
          type="submit"
          disabled={createRoom.isPending}
          className="rounded-lg bg-blue-600 px-5 py-2 font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        >
          สร้าง
        </button>
      </form>
      {errors.name && <p className="mb-3 text-sm text-red-500">{errors.name.message}</p>}

      {/* รายการห้อง */}
      {isLoading && <p className="text-slate-400">กำลังโหลด...</p>}
      {isError && <p className="text-red-500">โหลดห้องไม่สำเร็จ</p>}

      <ul className="space-y-2">
        {rooms?.map((room) => (
          <li key={room.id}>
            <Link
              to={`/rooms/${room.id}`}
              className="block rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-sm transition hover:border-blue-400 hover:shadow"
            >
              <span className="font-medium text-slate-800"># {room.name}</span>
            </Link>
          </li>
        ))}
        {rooms?.length === 0 && (
          <p className="text-slate-400">ยังไม่มีห้อง สร้างห้องแรกเลย!</p>
        )}
      </ul>
    </div>
  )
}
