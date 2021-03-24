package main

import (
	"context"

	pb "github.com/TRConley/clubhouse-backend-clone/gen/go"
)

// VoiceService implements the proto voice_service definition
type VoiceService struct {
}

// NewVoiceService will initialise a new voice service
func NewVoiceService() *VoiceService {
	return &VoiceService{}
}

// Hello will execute a simple request
func (v *VoiceService) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Out: in.In,
	}, nil
}
