package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	repositories "github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"

	"github.com/JanCalebManzano/tag-microservices/services/meta/proto"
)

type server struct {
	proto.UnimplementedMetaServiceServer

	repositoryClient repositories.MetaRepositoryClient
}

var _ proto.MetaServiceServer = (*server)(nil)

func (s server) GetAllSystems(ctx context.Context, _ *proto.GetAllSystemsRequest) (*proto.GetAllSystemsResponse, error) {
	result, err := s.repositoryClient.GetAllSystems(ctx, &repositories.GetAllSystemsRequest{})
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &proto.GetAllSystemsResponse{
		Systems: make([]*proto.System, len(result.Systems)),
	}

	for i, system := range result.Systems {
		res.Systems[i] = &proto.System{
			SystemNo:        system.SystemNo,
			SystemName:      system.SystemName,
			SystemShortName: system.SystemShortName,
			SetUser:         system.SetUser,
			SetTimestamp:    system.SetTimestamp,
		}
	}

	return res, nil
}
