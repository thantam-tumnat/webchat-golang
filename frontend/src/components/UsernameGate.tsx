import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useMutation } from '@tanstack/react-query'
import type { ReactNode } from 'react'
import { createUser } from '../api/endpoints'
import { useUserStore } from '../stores/userStore'

const schema = z.object({
  username: z
    .string()
    .min(2, 'ชื่อต้องยาวอย่างน้อย 2 ตัวอักษร')
    .max(50, 'ชื่อต้องไม่เกิน 50 ตัวอักษร'),
})
type FormValues = z.infer<typeof schema>

// UsernameGate: ถ้ายังไม่ได้ตั้ง username จะบังคับให้ตั้งก่อน
// ตั้งแล้วถึงจะแสดง children (ส่วนที่เหลือของแอป)
export function UsernameGate({ children }: { children: ReactNode }) {
  const { user, setUser } = useUserStore()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) })

  const mutation = useMutation({
    mutationFn: (username: string) => createUser(username),
    onSuccess: (u) => setUser(u),
  })

  // มี user แล้ว -> แสดงแอปได้เลย
  if (user) return <>{children}</>

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-100 p-4">
      <form
        onSubmit={handleSubmit((v) => mutation.mutate(v.username))}
        className="w-full max-w-sm rounded-2xl bg-white p-8 shadow-lg"
      >
        <h1 className="mb-2 text-2xl font-bold text-slate-800">💬 Chat App</h1>
        <p className="mb-6 text-sm text-slate-500">ตั้งชื่อที่จะใช้แสดงในแชท</p>

        <input
          {...register('username')}
          placeholder="username ของคุณ"
          autoFocus
          className="w-full rounded-lg border border-slate-300 px-4 py-2 outline-none focus:border-blue-500"
        />
        {errors.username && (
          <p className="mt-1 text-sm text-red-500">{errors.username.message}</p>
        )}

        <button
          type="submit"
          disabled={mutation.isPending}
          className="mt-4 w-full rounded-lg bg-blue-600 py-2 font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {mutation.isPending ? 'กำลังเข้า...' : 'เข้าใช้งาน'}
        </button>

        {mutation.isError && (
          <p className="mt-2 text-sm text-red-500">เกิดข้อผิดพลาด ลองใหม่อีกครั้ง</p>
        )}
      </form>
    </div>
  )
}
