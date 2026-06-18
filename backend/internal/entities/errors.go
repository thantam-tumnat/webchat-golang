package entities

// AppError คือ error มาตรฐานของระบบ มี HTTP status + code + message
// ทำให้ทุก error ที่ส่งกลับ client มีรูปแบบเดียวกัน { "code": ..., "message": ... }
type AppError struct {
	Status  int    `json:"-"`                 // HTTP status code เช่น 404, 400
	Code    string `json:"code"`              // รหัส error สำหรับให้ frontend เอาไปเช็ค เช่น "ROOM_NOT_FOUND"
	Message string `json:"message"`           // ข้อความอธิบายแบบ human-readable
	Details any    `json:"details,omitempty"` // รายละเอียดเพิ่มเติม เช่น list ของ field ที่ validate ไม่ผ่าน
}

func (e *AppError) Error() string {
	return e.Message
}

// WithDetails คืน AppError ตัวใหม่ (copy) ที่แนบ details เข้าไป
// ใช้ตอน validate ไม่ผ่าน เพื่อบอกว่า field ไหนผิดบ้าง
func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{
		Status:  e.Status,
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

// NewAppError สร้าง AppError ตัวใหม่
func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

// error สำเร็จรูปที่ใช้บ่อย
var (
	ErrUserNotFound  = NewAppError(404, "USER_NOT_FOUND", "user not found")
	ErrRoomNotFound  = NewAppError(404, "ROOM_NOT_FOUND", "room not found")
	ErrValidation    = NewAppError(400, "VALIDATION_ERROR", "validation failed")
	ErrInternal      = NewAppError(500, "INTERNAL_ERROR", "internal server error")
	ErrUsernameTaken = NewAppError(409, "USERNAME_TAKEN", "username already taken")
)
