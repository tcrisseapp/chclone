package repositories

import (
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/core/domain"
	"gorm.io/gorm"
)

// RoomRepository contains all the methods for the Room Repo
type RoomRepository struct {
	db *gorm.DB
}

// NewRoomRepository will initiate a new room service
func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

// Create will create a new room
func (r *RoomRepository) Create(room *domain.Room) (*domain.Room, error) {
	err := r.db.Create(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetBySID will retrieve the room by SID
func (r *RoomRepository) GetBySID(sid string) (*domain.Room, error) {
	var room *domain.Room
	err := r.db.Where("sid = ?", sid).First(room).Error
	if err != nil {
		return nil, err
	}
	return room, nil
}

// List will list all the rooms
func (r *RoomRepository) List() ([]*domain.Room, error) {
	var rooms []*domain.Room
	err := r.db.Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
