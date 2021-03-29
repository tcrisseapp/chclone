package main

import (
	"fmt"
	"os"

	"github.com/TRConley/clubhouse-backend-clone/cmd/room/handlers"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/models"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/repositories"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/services"
	pb "github.com/TRConley/clubhouse-backend-clone/gen/go/room/v1"
	"github.com/TRConley/clubhouse-backend-clone/pkg/database"
	"github.com/TRConley/clubhouse-backend-clone/pkg/grpcservice"
	"github.com/TRConley/clubhouse-backend-clone/pkg/zaplogger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// flags definition
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.StringP("config", "c", "config.yaml", "config file including the path")
	fs.String("level", "info", "log level debug, info, warn, error, flat or panic")

	// parse flags
	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	}

	// bind flags and environment variables
	viper.BindPFlags(fs)
	viper.AutomaticEnv()
	viper.SetConfigFile(viper.GetString("config"))
	if _, err := os.Stat(viper.GetString("config")); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config file: %s\n\n", err.Error())
		os.Exit(1)
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s\n\n", err.Error())
		os.Exit(1)
	}

	// configure the logger
	logger, err := zaplogger.New(viper.GetString("level"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %s\n\n", err.Error())
		os.Exit(1)
	}
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	logger.Info("starting service", zap.String("name", viper.GetString("svc")))

	// init the database
	db, err := initDatabase()
	if err != nil {
		logger.Panic("unable to init database", zap.Error(err))
	}

	roomRepository := repositories.NewRoomRepository(db)
	roomService := services.NewRoomService(roomRepository)

	// init the grpc service
	var grpcConfig grpcservice.Config
	if err := viper.Unmarshal(&grpcConfig); err != nil {
		logger.Fatal("error loading grpc config:", zap.Error(err))
	}

	grpcHandler := handlers.NewGRPCHandler(logger, roomService)
	s := grpcservice.NewService(&grpcConfig, logger)
	s.GRPCServer.RegisterService(&pb.RoomService_ServiceDesc, grpcHandler)

	s.ListenAndServe(grpcHandler)
}

func initDatabase() (*gorm.DB, error) {
	var databaseConfig database.Config
	if err := viper.Unmarshal(&databaseConfig); err != nil {
		return nil, err
	}

	db, err := database.Init(&databaseConfig)
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	db.AutoMigrate(&models.Room{})

	return db, nil
}
