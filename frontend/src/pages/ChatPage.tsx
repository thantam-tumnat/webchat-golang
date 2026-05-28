import { useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { useMessages, useSendMessage } from '../hooks/useMessages'
import { useUserStore } from '../stores/userStore'

export function ChatPage() {
  const { id } = useParams()
  const roomId = Number(id)
  const { user } = useUserStore()

  const { data, isLoading, isError } = useMessages(roomId)
  const sendMessage = useSendMessage(roomId)
  const [text, setText] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)

  // backend ส่งข้อความเรียงใหม่->เก่า เรากลับด้านให้เก่า->ใหม่ เพื่อแสดงแบบ chat ปกติ
  const messages = data ? [...data.data].reverse() : []

  // เลื่อนลงล่างสุดเมื่อมีข้อความใหม่
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages.length])

  function handleSend(e: React.FormEvent) {
    e.preventDefault()
    const content = text.trim()
    if (!content || !user) return
    sendMessage.mutate(
      { userId: user.id, content },
      { onSuccess: () => setText('') },
    )
  }

  return (
    <div className="mx-auto flex h-screen max-w-2xl flex-col p-4">
      <header className="mb-3 flex items-center gap-3">
        <Link to="/" className="text-slate-500 hover:text-slate-800">
          ← กลับ
        </Link>
        <h1 className="text-lg font-bold text-slate-800">ห้อง #{roomId}</h1>
      </header>

      {/* กล่องข้อความ */}
      <div className="flex-1 space-y-2 overflow-y-auto rounded-xl border border-slate-200 bg-white p-4">
        {isLoading && <p className="text-slate-400">กำลังโหลด...</p>}
        {isError && <p className="text-red-500">โหลดข้อความไม่สำเร็จ</p>}

        {messages.map((m) => {
          const isMine = m.user_id === user?.id
          return (
            <div
              key={m.id}
              className={`flex flex-col ${isMine ? 'items-end' : 'items-start'}`}
            >
              <span className="text-xs text-slate-400">
                {m.user?.username ?? `user ${m.user_id}`}
              </span>
              <span
                className={`max-w-[75%] rounded-2xl px-4 py-2 ${
                  isMine
                    ? 'bg-blue-600 text-white'
                    : 'bg-slate-100 text-slate-800'
                }`}
              >
                {m.content}
              </span>
            </div>
          )
        })}
        <div ref={bottomRef} />
      </div>

      {/* ฟอร์มส่งข้อความ */}
      <form onSubmit={handleSend} className="mt-3 flex gap-2">
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="พิมพ์ข้อความ..."
          className="flex-1 rounded-lg border border-slate-300 px-4 py-2 outline-none focus:border-blue-500"
        />
        <button
          type="submit"
          disabled={sendMessage.isPending}
          className="rounded-lg bg-blue-600 px-5 py-2 font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        >
          ส่ง
        </button>
      </form>
    </div>
  )
}
