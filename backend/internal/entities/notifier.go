package entities

// MessageNotifier เป็น "สัญญา" ให้ usecase แจ้งเตือนเมื่อมีข้อความใหม่
// usecase รู้จักแค่ interface นี้ ไม่รู้ว่าเบื้องหลังส่งผ่านอะไร (WebSocket, SSE, ฯลฯ)
//
// วงนอก (delivery/websocket) เป็นคนมา implement แล้ว main เป็นคนเชื่อมให้

type MessageNotifier interface {
	// NotifyMessage ถูกเรียกหลังบันทึกข้อความสำเร็จ เพื่อกระจายให้ทุกคนในห้องนั้นเห็นทันที
	NotifyMessage(roomID uint, msg *Message)
}
