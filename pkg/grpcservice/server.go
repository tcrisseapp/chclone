package grpcservice

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/TRConley/clubhouse-backend-clone/pkg/signals"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Config specifies the
type Config struct {
	ServiceName string `mapstructure:"svc"`
	GRPCPort    int    `mapstructure:"grpc-port"`
	HTTPPort    int    `mapstructure:"http-port"`
}

// Registry holds the reference to services that implement protot
type Registry interface {
	Config() *Config
	RegisterGRPC(*grpc.Server)
	RegisterHTTP(ctx context.Context, gwMux *runtime.ServeMux, conn *grpc.ClientConn)
}

// ListenAndServe will listen and serve the GRPC registry
func ListenAndServe(registry Registry, logger *zap.Logger) {
	stopCh := signals.SetupSignalHandler()
	config := registry.Config()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Int("port", config.GRPCPort))
	}

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	reflection.Register(grpcServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	registry.RegisterGRPC(grpcServer)

	// Start serving the gRPC server
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	logger.Info("started gRPC server", zap.Int("port", config.GRPCPort))

	// Connect to the gRPC server
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", config.GRPCPort),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Fatal("failed to dial server", zap.Error(err))
	}

	gwmux := runtime.NewServeMux()
	registry.RegisterHTTP(context.Background(), gwmux, conn)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTPPort),
		Handler: gwmux,
	}

	// Start serving the HTTP server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatal("failed to serve http server", zap.Error(err))
		}
	}()

	logger.Info("started HTTP server", zap.Int("port", config.HTTPPort))

	<-stopCh
	grpcServer.GracefulStop()
	httpServer.Shutdown(context.Background())

	logger.Info("shutting down gRPC and HTTP server")
}
