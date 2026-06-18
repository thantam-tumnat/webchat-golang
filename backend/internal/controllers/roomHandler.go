package controllers

import (
	"chatapp/internal/entities"
	"chatapp/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	roomUsecase *usecases.RoomUsecase
}

func NewRoomHandler(roomUsecase *usecases.RoomUsecase) *RoomHandler {
	return &RoomHandler{roomUsecase: roomUsecase}
}

type createRoomRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// List คืนรายการห้องทั้งหมด
// GET /api/rooms
func (h *RoomHandler) List(c *fiber.Ctx) error {
	rooms, err := h.roomUsecase.List(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

// Create สร้างห้องใหม่
// POST /api/rooms
func (h *RoomHandler) Create(c *fiber.Ctx) error {
	var req createRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return entities.ErrValidation
	}
	if err := validateStruct(req); err != nil {
		return err
	}

	room, err := h.roomUsecase.Create(c.Context(), req.Name)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(room)
}
