package grpc

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials/insecure"

	repositories "github.com/JanCalebManzano/tag-microservices/repositories/meta/proto"

	pkggrpc "github.com/JanCalebManzano/tag-microservices/pkg/grpc"
	"github.com/JanCalebManzano/tag-microservices/services/meta/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, port int, logger *zap.Logger) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	conn, err := grpc.DialContext(ctx,
		fmt.Sprintf("localhost:%s", os.Getenv("APPS_META_REPOSITORY_PORT")), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial db server: %w", err)
	}

	repositoryClient := repositories.NewMetaRepositoryClient(conn)

	svc := &server{
		repositoryClient: repositoryClient,
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterMetaServiceServer(s, svc)
	}).Start(ctx)
}
