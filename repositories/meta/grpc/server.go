package grpc

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/db"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"
)

type server struct {
	proto.UnimplementedMetaRepositoryServer

	log                 *zap.Logger
	db                  db.MetaDB
	systemSubscriptions map[proto.MetaRepository_SubscribeSystemsServer][]*proto.GetAllSystemsRequest
}

var _ proto.MetaRepositoryServer = (*server)(nil)

func newServer(ctx context.Context, log *zap.Logger, metaDB db.MetaDB) *server {
	s := &server{
		log:                 log,
		db:                  metaDB,
		systemSubscriptions: make(map[proto.MetaRepository_SubscribeSystemsServer][]*proto.GetAllSystemsRequest, 0),
	}
	go s.handleUpdates(ctx)

	return s
}

func (s *server) handleUpdates(ctx context.Context) {
	for range s.db.MonitorSystems(5 * time.Second) {
		s.log.Info("server: Received updated systems")

		for k, v := range s.systemSubscriptions {
			for _, srs := range v {
				systems, err := s.GetAllSystems(ctx, srs)
				if err != nil {
					s.log.Error("server: Unable to get updated systems", zap.Error(err))
				}

				if err := k.Send(systems); err != nil {
					s.log.Error("server: Unable to set updated systems")
				}
			}
		}
	}
}

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

func (s *server) SubscribeSystems(stream proto.MetaRepository_SubscribeSystemsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.log.Info("Client has closed connection")
			delete(s.systemSubscriptions, stream)
			break
		}

		if err != nil {
			s.log.Error("Unable to read from client", zap.Error(err))
			delete(s.systemSubscriptions, stream)
			return err
		}

		s.log.Info("Handle client request", zap.Any("request", req))

		reqs, ok := s.systemSubscriptions[stream]
		if !ok {
			reqs = make([]*proto.GetAllSystemsRequest, 0)
		}

		reqs = append(reqs, req)
		s.systemSubscriptions[stream] = reqs
	}

	return nil
}
