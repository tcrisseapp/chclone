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

// HTTPGateway specifies the method in order to implement HTTP gateway functionality
type HTTPGateway interface {
	RegisterHTTP(ctx context.Context, gwMux *runtime.ServeMux, conn *grpc.ClientConn)
}

// Service holds to properties for the GRPC service
type Service struct {
	config     *Config
	logger     *zap.Logger
	GRPCServer *grpc.Server
}

// NewService will initiate a new GRPC service
func NewService(config *Config, logger *zap.Logger) *Service {
	return &Service{
		config:     config,
		logger:     logger,
		GRPCServer: grpc.NewServer(),
	}
}

// ListenAndServe will listen and serve the GRPC registry
func (s *Service) ListenAndServe(httpService HTTPGateway) {
	stopCh := signals.SetupSignalHandler()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GRPCPort))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Int("port", s.config.GRPCPort))
	}

	healthServer := health.NewServer()
	reflection.Register(s.GRPCServer)
	grpc_health_v1.RegisterHealthServer(s.GRPCServer, healthServer)
	healthServer.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	// Start serving the gRPC server
	go func() {
		if err := s.GRPCServer.Serve(listener); err != nil {
			s.logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	s.logger.Info("started gRPC server", zap.Int("port", s.config.GRPCPort))

	// Connect to the gRPC server
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", s.config.GRPCPort),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		s.logger.Fatal("failed to dial server", zap.Error(err))
	}

	gwmux := runtime.NewServeMux()
	httpService.RegisterHTTP(context.Background(), gwmux, conn)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.HTTPPort),
		Handler: gwmux,
	}

	// Start serving the HTTP server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			s.logger.Fatal("failed to serve http server", zap.Error(err))
		}
	}()

	s.logger.Info("started HTTP server", zap.Int("port", s.config.HTTPPort))

	<-stopCh
	s.GRPCServer.GracefulStop()
	httpServer.Shutdown(context.Background())

	s.logger.Info("shutting down gRPC and HTTP server")
}
