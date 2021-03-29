package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Room holds the properties for a Room entity
type Room struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
