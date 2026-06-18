package repository

import (
	"context"
	"errors"

	"chatapp/internal/domain"

	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepository) FindAll(ctx context.Context) ([]domain.Room, error) {
	var rooms []domain.Room
	// เรียงห้องใหม่สุดไว้บน
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) FindByID(ctx context.Context, id uint) (*domain.Room, error) {
	var room domain.Room
	err := r.db.WithContext(ctx).First(&room, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRoomNotFound
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}
