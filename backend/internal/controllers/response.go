package http

import (
	"errors"

	"chatapp/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// validate เป็น instance กลางของ validator ใช้ตรวจ struct request
var validate = validator.New()

// ErrorResponse คือรูปแบบ error เดียวกันทั้งระบบ
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// PaginatedResponse รูปแบบ response สำหรับข้อมูลที่แบ่งหน้า
type PaginatedResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// ErrorHandler เป็น custom error handler กลางของ Fiber
// ทุก error ที่ handler return จะวิ่งมาที่นี่ แล้วถูกแปลงเป็น JSON รูปแบบเดียวกัน
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 1) error ของระบบเราเอง (AppError)
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Status).JSON(ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		})
	}

	// 2) error ของ Fiber เอง (เช่น 404 route not found)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(ErrorResponse{
			Code:    "HTTP_ERROR",
			Message: fiberErr.Message,
		})
	}

	// 3) error อื่นๆ ที่ไม่คาดคิด → 500 (ไม่เปิดเผยรายละเอียดภายในให้ client)
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Code:    domain.ErrInternal.Code,
		Message: domain.ErrInternal.Message,
	})
}

// validateStruct ตรวจ struct ด้วย validator ถ้าไม่ผ่านคืน AppError 400 พร้อม details
func validateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		var fields []string
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				fields = append(fields, fe.Field()+" ("+fe.Tag()+")")
			}
		}
		return domain.ErrValidation.WithDetails(fields)
	}
	return nil
}
