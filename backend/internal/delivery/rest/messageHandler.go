package rest

import (
	"chatapp/internal/entities"
	"chatapp/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	messageUsecase *usecases.MessageUsecase
}

func NewMessageHandler(messageUsecase *usecases.MessageUsecase) *MessageHandler {
	return &MessageHandler{messageUsecase: messageUsecase}
}

type sendMessageRequest struct {
	UserID  uint   `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=2000"`
}

// List ดึงข้อความในห้องแบบแบ่งหน้า
// GET /api/rooms/:id/messages?page=1&limit=20
func (h *MessageHandler) List(c *fiber.Ctx) error {
	roomID, err := c.ParamsInt("id")
	if err != nil || roomID <= 0 {
		return entities.ErrRoomNotFound
	}

	// อ่าน query param page/limit พร้อม default + กันค่าผิด
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 20)
	if limit < 1 || limit > 100 {
		limit = 20
	}

	messages, total, err := h.messageUsecase.List(c.Context(), uint(roomID), page, limit)
	if err != nil {
		return err
	}

	return c.JSON(PaginatedResponse{
		Data:  messages,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// Send ส่งข้อความใหม่เข้าห้อง
// POST /api/rooms/:id/messages
func (h *MessageHandler) Send(c *fiber.Ctx) error {
	roomID, err := c.ParamsInt("id")
	if err != nil || roomID <= 0 {
		return entities.ErrRoomNotFound
	}

	var req sendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return entities.ErrValidation
	}
	if err := validateStruct(req); err != nil {
		return err
	}

	msg, err := h.messageUsecase.Send(c.Context(), uint(roomID), req.UserID, req.Content)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(msg)
}
