package services

import (
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/domain"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/ports"
)

// RoomService will act as the entrypoint to the core domain
type RoomService struct {
	roomRepository ports.RoomRepository
}

// NewRoomService will initiate a new room service
func NewRoomService(roomRepo ports.RoomRepository) *RoomService {
	return &RoomService{
		roomRepository: roomRepo,
	}
}

// Create will create a new room
func (r *RoomService) Create(room *domain.Room) (*domain.Room, error) {
	return r.roomRepository.Create(room)
}

// List will list all the rooms
func (r *RoomService) List() ([]*domain.Room, error) {
	return r.roomRepository.List()
}
