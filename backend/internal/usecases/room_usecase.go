package usecase

import (
	"context"

	"chatapp/internal/domain"
)

type RoomUsecase struct {
	roomRepo domain.RoomRepository
}

func NewRoomUsecase(roomRepo domain.RoomRepository) *RoomUsecase {
	return &RoomUsecase{roomRepo: roomRepo}
}

func (uc *RoomUsecase) List(ctx context.Context) ([]domain.Room, error) {
	return uc.roomRepo.FindAll(ctx)
}

func (uc *RoomUsecase) Create(ctx context.Context, name string) (*domain.Room, error) {
	room := &domain.Room{Name: name}
	if err := uc.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}
