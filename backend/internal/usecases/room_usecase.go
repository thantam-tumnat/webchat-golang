package usecases

import (
	"context"

	"chatapp/internal/entities"
)

type RoomUsecase struct {
	roomRepo entities.RoomRepository
}

func NewRoomUsecase(roomRepo entities.RoomRepository) *RoomUsecase {
	return &RoomUsecase{roomRepo: roomRepo}
}

func (uc *RoomUsecase) List(ctx context.Context) ([]entities.Room, error) {
	return uc.roomRepo.FindAll(ctx)
}

func (uc *RoomUsecase) Create(ctx context.Context, name string) (*entities.Room, error) {
	room := &entities.Room{Name: name}
	if err := uc.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}
