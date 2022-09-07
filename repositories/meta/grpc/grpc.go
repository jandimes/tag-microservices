package grpc

import (
	"context"
	"fmt"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/db"

	pkggrpc "github.com/JanCalebManzano/tag-microservices/pkg/grpc"
	"github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger *zap.Logger) error {
	metaDB, err := db.NewMetaDB(ctx, logger)
	if err != nil {
		return fmt.Errorf("failed to connect to db server: %w", err)
	}

	svc := newServer(ctx, logger, metaDB)

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterMetaRepositoryServer(s, svc)
	}).Start(ctx)
}
