package handlers

import (
	"context"
	"io"

	"github.com/TRConley/clubhouse-backend-clone/cmd/room/core/domain"
	"github.com/TRConley/clubhouse-backend-clone/cmd/room/core/ports"
	pb "github.com/TRConley/clubhouse-backend-clone/gen/go/room/v1"
	sfuPB "github.com/TRConley/clubhouse-backend-clone/gen/go/sfu/v1"
	"github.com/TRConley/clubhouse-backend-clone/pkg/grpcservice"
	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// GRPCHandler will handle the gRPC input
type GRPCHandler struct {
	pb.UnimplementedRoomServiceServer

	logger      *zap.Logger
	roomService ports.RoomService
	sfuClient   sfuPB.SFUClient
}

// NewGRPCHandler will initiate a new GRPCHandler instance
func NewGRPCHandler(logger *zap.Logger, roomService ports.RoomService) (*GRPCHandler, error) {
	sfuConn, err := grpcservice.Dial("ion-sfu:50052")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to ion-suf")
	}

	return &GRPCHandler{
		logger:      logger,
		roomService: roomService,
		sfuClient:   sfuPB.NewSFUClient(sfuConn),
	}, nil
}

// Signal will function as a wrapper around the ion-sfu signal logic
func (g *GRPCHandler) Signal(stream pb.RoomService_SignalServer) error {
	g.logger.Info("came heree")

	ctx := context.Background()
	errorCh := make(chan error)
	replyCh := make(chan *sfuPB.SignalReply)
	requestCh := make(chan *sfuPB.SignalRequest)

	sfuStream, err := g.sfuClient.Signal(ctx)
	if err != nil {
		g.logger.Error("unable to open signal stream to sfu", zap.Error(err))
		return status.Error(codes.Internal, "unable to signal sfu")
	}

	defer func() {
		close(errorCh)
		close(replyCh)
		close(requestCh)
	}()

	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				g.logger.Error("failed to receive message", zap.Error(err))
				errorCh <- err
				return
			}
			requestCh <- req
		}
	}()

	go func() {
		for {
			rep, err := sfuStream.Recv()
			if err != nil {
				g.logger.Error("failed to reply message", zap.Error(err))
				errorCh <- err
				return
			}
			replyCh <- rep
		}
	}()

	for {
		select {
		case err := <-errorCh:
			g.logger.Error("recevied an error", zap.Error(err))
			return status.Error(codes.Internal, "unable to create room")
		case reply, ok := <-replyCh:
			if !ok {
				return io.EOF
			}
			stream.Send(reply)
		case request, ok := <-requestCh:
			if !ok {
				return io.EOF
			}
			g.logger.Info("received signal request", zap.Any("request", request.String()))

			// handle the request message
			switch requestM := request.Payload.(type) {
			case *sfuPB.SignalRequest_Join:
				g.logger.Info("received join request")

				_, err := g.roomService.GetBySID(requestM.Join.Sid)
				if err != nil && err != gorm.ErrRecordNotFound {
					g.logger.Error("unabel to get room by sid", zap.Error(err))
					return status.Error(codes.Internal, "unable to get room by sid")
				}

				if !errors.Is(err, gorm.ErrRecordNotFound) {
					_, err := g.roomService.Create(requestM.Join.Sid)
					if err != nil {
						g.logger.Error("unable to create room", zap.Error(err))
						return status.Error(codes.Internal, "unable to create room")
					}
				}

				err = sfuStream.Send(request)
				if err != nil {
					g.logger.Error("unable to send message to sfu", zap.Error(err))
					return status.Error(codes.Internal, "unable to forward message to sfu")
				}
			default:
				g.logger.Info("received other request")

				err := sfuStream.Send(request)
				if err != nil {
					g.logger.Error("unable to send message to sfu", zap.Error(err))
					return status.Error(codes.Internal, "unable to forward message to sfu")
				}

			}
		}
	}

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
