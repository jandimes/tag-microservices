package grpc

import (
	"context"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	repositories "github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"

	"github.com/JanCalebManzano/tag-microservices/services/meta/proto"
)

type server struct {
	proto.UnimplementedMetaServiceServer

	log     *zap.Logger
	repo    repositories.MetaRepositoryClient
	systems []*repositories.System
}

var _ proto.MetaServiceServer = (*server)(nil)

func newServer(ctx context.Context, log *zap.Logger, repo repositories.MetaRepositoryClient) *server {
	s := &server{
		log:     log,
		repo:    repo,
		systems: make([]*repositories.System, 0),
	}
	go s.handleUpdates(ctx)

	return s
}

func (s *server) handleUpdates(ctx context.Context) {
	stream, err := s.repo.SubscribeSystems(ctx)
	if err != nil {
		s.log.Error("client: Unable to subscribe to systems", zap.Error(err))
		return
	}

	if err := stream.Send(&repositories.GetAllSystemsRequest{}); err != nil {
		s.log.Error("client: Unable to set updated systems")
		return
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			s.log.Error("client: Error receiving systems", zap.Error(err))
			return
		}

		s.log.Info("client: Received updated systems")
		s.systems = req.Systems
	}
}

func (s *server) GetAllSystems(
	ctx context.Context, _ *proto.GetAllSystemsRequest) (*proto.GetAllSystemsResponse, error) {
	systems := make([]*repositories.System, 0)

	if len(s.systems) != 0 {
		systems = s.systems
	} else {
		result, err := s.repo.GetAllSystems(ctx, &repositories.GetAllSystemsRequest{})
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}

		if len(result.Systems) == 0 {
			return &proto.GetAllSystemsResponse{
				Status: "failure",
				Errors: []*proto.ResponseError{
					{
						StatusCode:   "E00002",
						ErrorMessage: "No data present. ",
					},
				},
			}, nil
		}

		systems = result.Systems
	}

	res := &proto.GetAllSystemsResponse{
		Status: "success",
		Data:   make([]*proto.System, len(systems)),
	}

	for i, system := range systems {
		res.Data[i] = &proto.System{
			SystemNo:        system.SystemNo,
			SystemName:      system.SystemName,
			SystemShortName: system.SystemShortName,
			SetUser:         system.SetUser,
			SetTimestamp:    system.SetTimestamp,
		}
	}

	return res, nil
}
