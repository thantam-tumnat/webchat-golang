package rest

import (
	"chatapp/internal/entities"
	"chatapp/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

// UserHandler รับ HTTP request เกี่ยวกับ user แล้วเรียก usecase
type UserHandler struct {
	userUsecase *usecases.UserUsecase
}

func NewUserHandler(userUsecase *usecases.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// createUserRequest คือ body ที่รับมาจาก POST /api/users
type createUserRequest struct {
	Username string `json:"username" validate:"required,min=2,max=50"`
}

// Create สร้าง user ใหม่ (หรือคืนตัวเดิมถ้า username ซ้ำ)
// POST /api/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return entities.ErrValidation
	}
	if err := validateStruct(req); err != nil {
		return err
	}

	user, err := h.userUsecase.CreateOrGet(c.Context(), req.Username)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}
