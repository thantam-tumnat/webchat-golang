package repositories

import (
	"context"
	"errors"

	"chatapp/internal/entities"

	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) entities.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *entities.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepository) FindAll(ctx context.Context) ([]entities.Room, error) {
	var rooms []entities.Room
	// เรียงห้องใหม่สุดไว้บน
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) FindByID(ctx context.Context, id uint) (*entities.Room, error) {
	var room entities.Room
	err := r.db.WithContext(ctx).First(&room, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entities.ErrRoomNotFound
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}
