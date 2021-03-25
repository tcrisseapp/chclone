package main

import (
	"context"
	"errors"

	pb "github.com/TRConley/clubhouse-backend-clone/gen/go/backend/v1"
	"github.com/TRConley/clubhouse-backend-clone/pkg/database"
	"github.com/TRConley/clubhouse-backend-clone/svc/backend/models"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
	ErrInternal    = errors.New("internal server error")
)

// BackendService implements the proto backend definition
type BackendService struct {
	pb.UnimplementedBackendServiceServer

	db     *gorm.DB
	logger *zap.Logger
}

// NewBackendService will initialise a new voice service
func NewBackendService(db *gorm.DB, logger *zap.Logger) *BackendService {
	return &BackendService{
		db:     db,
		logger: logger,
	}
}

// CreateRoom will create a new audio room
func (b *BackendService) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, ErrNameIsEmpty.Error())
	}

	room := &models.Room{
		BaseModel: database.BaseModel{
			ID: uuid.Must(uuid.NewV4(), nil),
		},
		Name: in.GetName(),
	}

	if err := models.AddRoom(b.db, room); err != nil {
		b.logger.Error("failed to create room", zap.Error(err))
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	roomPb, err := convertRoomToPB(room)
	if err != nil {
		b.logger.Error("failed to convert room to proto", zap.Error(err))
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	return &pb.CreateRoomResponse{Room: roomPb}, nil
}

// ListRooms will list all the available rooms
func (b *BackendService) ListRooms(ctx context.Context, in *empty.Empty) (*pb.ListRoomsResponse, error) {
	rooms, err := models.ListRooms(b.db)
	if err != nil {
		b.logger.Error("failed to list rooms", zap.Error(err))
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	var out pb.ListRoomsResponse
	for i := range rooms {
		outRoom, err := convertRoomToPB(rooms[i])
		if err != nil {
			b.logger.Error("failed to convert room to proto", zap.Error(err))
			return nil, status.Error(codes.Internal, ErrInternal.Error())

		}
		out.Rooms = append(out.Rooms, outRoom)
	}

	return &out, nil
}

func convertRoomToPB(room *models.Room) (*pb.Room, error) {
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
