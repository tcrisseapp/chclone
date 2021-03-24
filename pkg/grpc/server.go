package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Config contains the properties to configure the GRPC server
type Config struct {
	Port        int
	ServiceName string
}

// Server contains properties to start a GRPC server
type Server struct {
	*grpc.Server
	logger *zap.Logger
	config *Config
}

// NewServer will initiate a new Server instance
func NewServer(config *Config, logger *zap.Logger) *Server {
	return &Server{
		Server: grpc.NewServer(),
		logger: logger,
		config: config,
	}
}

// ListenAndServe will listen and serve the GRPC server
func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Int("port", s.config.Port))
	}

	healthSrv := health.NewServer()
	reflection.Register(s.Server)
	grpc_health_v1.RegisterHealthServer(s, healthSrv)
	healthSrv.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		if err := s.Serve(listener); err != nil {
			s.logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	<-stopCh
	s.Server.GracefulStop()

	s.logger.Info("shutting down gRPC server")
}
