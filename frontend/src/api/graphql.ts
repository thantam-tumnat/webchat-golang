// GraphQL client เล็ก ๆ ด้วย fetch — ไม่ต้องลง Apollo/urql
// ยิง POST ไป /graphql (vite proxy ส่งต่อไป backend) พร้อม query + variables
// แล้วคืนเฉพาะส่วน data; ถ้า GraphQL ตอบ errors มาก็ throw ออกไปให้ React Query จัดการ
export async function gql<T>(
  query: string,
  variables?: Record<string, unknown>,
): Promise<T> {
  const res = await fetch('/graphql', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables }),
  })

  const json = await res.json()
  if (json.errors?.length) {
    throw new Error(json.errors[0]?.message ?? 'GraphQL error')
  }
  return json.data as T
}
