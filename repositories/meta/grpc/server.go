package grpc

import (
	"context"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"
)

type server struct {
	proto.UnimplementedMetaRepositoryServer

	db db.MetaDB
}

var _ proto.MetaRepositoryServer = (*server)(nil)

func (s *server) GetAllSystems(ctx context.Context, _ *proto.GetAllSystemsRequest) (*proto.GetAllSystemsResponse, error) {
	systems, err := s.db.GetAllSystems(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &proto.GetAllSystemsResponse{
		Systems: make([]*proto.System, len(systems)),
	}

	for i, system := range systems {
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
