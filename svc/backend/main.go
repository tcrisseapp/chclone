package main

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	pb "github.com/TRConley/clubhouse-backend-clone/gen/go/backend/v1"
	"github.com/TRConley/clubhouse-backend-clone/svc/backend/models"

	"github.com/TRConley/clubhouse-backend-clone/pkg/database"
	"github.com/TRConley/clubhouse-backend-clone/pkg/grpcservice"
	"github.com/TRConley/clubhouse-backend-clone/pkg/zaplogger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	// init the configuration
	initConfig()

	// configure the logger
	logger, err := zaplogger.New(viper.GetString("level"))
	if err != nil {
		panic(errors.Wrap(err, "unable to init zap"))
	}
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	// init the database
	db, err := initDatabase()
	if err != nil {
		logger.Panic("unable to init database", zap.Error(err))
	}

	// init the gRPC regisry
	registry, err := initRegistry(db, logger)
	if err != nil {
		logger.Panic("unable to init registry", zap.Error(err))
	}

	// starting the gRPC server
	grpcservice.ListenAndServe(registry, logger)
}

// BackendRegistry contains the proto implementation
type BackendRegistry struct {
	config         *grpcservice.Config
	backendService *BackendService
}

// Config will return the grpc configuration
func (v *BackendRegistry) Config() *grpcservice.Config {
	return v.config
}

// RegisterGRPC will register the proto implementation on given gRPC server
func (v *BackendRegistry) RegisterGRPC(grpcServer *grpc.Server) {
	pb.RegisterBackendServiceServer(grpcServer, v.backendService)
}

// RegisterHTTP will register the HTTP server for service
func (v *BackendRegistry) RegisterHTTP(ctx context.Context, gwMux *runtime.ServeMux, conn *grpc.ClientConn) {
	pb.RegisterBackendServiceHandler(ctx, gwMux, conn)
}

func initRegistry(db *gorm.DB, logger *zap.Logger) (*BackendRegistry, error) {
	var grpcConfig grpcservice.Config
	if err := viper.Unmarshal(&grpcConfig); err != nil {
		return nil, err
	}

	registry := &BackendRegistry{
		config:         &grpcConfig,
		backendService: NewBackendService(db, logger),
	}

	return registry, nil
}

func initConfig() {
	viper.SetConfigFile("config.yaml")
	viper.SetDefault("svc", "placeholder")
	viper.SetDefault("level", "info")
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "unable to init config"))
	}
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
