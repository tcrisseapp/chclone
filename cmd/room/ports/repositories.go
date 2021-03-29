package ports

import "github.com/TRConley/clubhouse-backend-clone/cmd/room/domain"

// RoomRepository specifies the room repository operations
type RoomRepository interface {
	Create(room *domain.Room) (*domain.Room, error)
	List() ([]*domain.Room, error)
}
