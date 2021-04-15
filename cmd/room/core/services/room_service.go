package services

import (
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/core/domain"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/core/ports"
	uuid "github.com/satori/go.uuid"
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
func (r *RoomService) Create(sid string) (*domain.Room, error) {
	room := &domain.Room{
		ID:   uuid.Must(uuid.NewV4(), nil),
		SID:  sid,
		Name: "something",
	}
	return r.roomRepository.Create(room)
}

// GetBySID will retrieve the room with specified SID
func (r *RoomService) GetBySID(sid string) (*domain.Room, error) {
	return r.roomRepository.GetBySID(sid)
}

// List will list all the rooms
func (r *RoomService) List() ([]*domain.Room, error) {
	return r.roomRepository.List()
}
