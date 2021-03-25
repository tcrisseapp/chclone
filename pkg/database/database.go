package database

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connString used to inititate the postgres db connection
var connString = "host=%v port=%d user=%v dbname=%v password=%v"

// Config holds the configuration for the database
type Config struct {
	Hostname string `mapstructure:"pg-host"`
	Username string `mapstructure:"pg-user"`
	Password string `mapstructure:"pg-pass"`
	DBName   string `mapstructure:"pg-dbname"`
	Port     int    `mapstructure:"pg-port"`
}

// BaseModel holds the base properties needed for each model
type BaseModel struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Init will initialise a new database
func Init(config *Config) (*gorm.DB, error) {
	conn := fmt.Sprintf(connString, config.Hostname, config.Port, config.Username, config.DBName, config.Password)
	fmt.Println(conn)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  conn,
		PreferSimpleProtocol: true,
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}
