package grpcservice

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/TRConley/clubhouse-backend-clone/pkg/signals"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Config specifies the
type Config struct {
	ServiceName string `mapstructure:"svc"`
	Port        int    `mapstructure:"port"`
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
func (s *Service) ListenAndServe() {
	stopCh := signals.SetupSignalHandler()

	healthServer := health.NewServer()
	reflection.Register(s.GRPCServer)
	grpc_health_v1.RegisterHealthServer(s.GRPCServer, healthServer)
	healthServer.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	options := []grpcweb.Option{
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithWebsockets(true),
	}

	wrappedServer := grpcweb.WrapServer(s.GRPCServer, options...)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.config.Port),
		Handler: http.HandlerFunc(handler),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Int("port", s.config.Port))
	}

	m := cmux.New(listener)

	// Start serving the gRPC server
	grpcListener := m.Match(cmux.HTTP2())
	go func() {
		if err := s.GRPCServer.Serve(grpcListener); err != nil {
			s.logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	// Start serving th HTTP server
	httpListener := m.Match(cmux.HTTP1Fast())
	go func() {
		if err := httpServer.Serve(httpListener); err != nil {
			s.logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	s.logger.Info("started gRPC server", zap.Int("port", s.config.Port))

	<-stopCh
	s.GRPCServer.GracefulStop()
	httpServer.Shutdown(context.Background())
	s.logger.Info("shutting down gRPC server")
}
