package handlers

import (
	"context"

	"github.com/TRConley/clubhouse-backend-clone/cmd/room/domain"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/ports"
	pb "github.com/TRConley/clubhouse-backend-clone/gen/go/room/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler will handle the gRPC input
type GRPCHandler struct {
	pb.UnimplementedRoomServiceServer

	logger      *zap.Logger
	roomService ports.RoomService
}

// NewGRPCHandler will initiate a new GRPCHandler instance
func NewGRPCHandler(logger *zap.Logger, roomService ports.RoomService) *GRPCHandler {
	return &GRPCHandler{
		logger:      logger,
		roomService: roomService,
	}
}

// RegisterHTTP will register the HTTP grpc-gateway
func (g *GRPCHandler) RegisterHTTP(ctx context.Context, gwMux *runtime.ServeMux, conn *grpc.ClientConn) {
	pb.RegisterRoomServiceHandler(ctx, gwMux, conn)
}

// CreateRoom will create a new room
func (g *GRPCHandler) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}

	room, err := g.roomService.Create(&domain.Room{
		ID:   uuid.Must(uuid.NewV4(), nil),
		Name: in.GetName(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create room")
	}

	out, err := convertRoomToPB(room)
	if err != nil {
		g.logger.Error("failed to convert room to proto", zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to create room")
	}

	return &pb.CreateRoomResponse{
		Room: out,
	}, nil
}

// ListRooms will list all the rooms
func (g *GRPCHandler) ListRooms(ctx context.Context, in *empty.Empty) (*pb.ListRoomsResponse, error) {
	rooms, err := g.roomService.List()
	if err != nil {
		g.logger.Error("failed to list rooms", zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to list rooms")
	}

	var out pb.ListRoomsResponse
	for i := range rooms {
		outRoom, err := convertRoomToPB(rooms[i])
		if err != nil {
			g.logger.Error("failed to convert room to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, "unable to list rooms")

		}
		out.Rooms = append(out.Rooms, outRoom)
	}

	return &out, nil
}

func convertRoomToPB(room *domain.Room) (*pb.Room, error) {
	pCreatedAt, err := ptypes.TimestampProto(room.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &pb.Room{
		Id:        room.ID.String(),
		Name:      room.Name,
		CreatedAt: pCreatedAt,
	}, nil
}
