package grpc

import (
	"context"

	pkggrpc "github.com/JanCalebManzano/tag-microservices/pkg/grpc"
	"github.com/JanCalebManzano/tag-microservices/repositories/meta/db"
	"github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger *zap.Logger) error {
	metaDB, err := db.NewMetaDB()
	if err != nil {
		return err
	}

	svc := &server{
		db: metaDB,
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterMetaRepositoryServer(s, svc)
	}).Start(ctx)
}
