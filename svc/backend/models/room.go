package models

import (
	"github.com/TRConley/clubhouse-backend-clone/pkg/database"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Room holds the properties for a Room entity
type Room struct {
	database.BaseModel
	Name string
}

// CreateRoom will add a new room
func AddRoom(db *gorm.DB, room *Room) error {
	return db.Create(room).Error
}

// GetRoom will get room by specified ID
func GetRoom(db *gorm.DB, id uuid.UUID) (*Room, error) {
	var room Room
	err := db.First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// ListRooms will list all the roomrs
func ListRooms(db *gorm.DB) ([]*Room, error) {
	var rooms []*Room
	err := db.Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
