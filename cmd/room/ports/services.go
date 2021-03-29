// Package services specifies all the services within the Room service
package ports

import "github.com/TRConley/clubhouse-backend-clone/cmd/room/domain"

// RoomService specifies the room service operations
type RoomService interface {
	Create(room *domain.Room) (*domain.Room, error)
	List() ([]*domain.Room, error)
}
