package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	metapb "github.com/JanCalebManzano/tag-microservices/services/meta/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func RunServer(ctx context.Context, port int, _ *zap.Logger) error {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: false,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	metaConn, err := grpc.DialContext(ctx,
		fmt.Sprintf("localhost:%s", os.Getenv("APPS_META_SERVICE_PORT")), opts...)
	if err != nil {
		return fmt.Errorf("failed to dial to meta grpc server: %w", err)
	}

	if err := metapb.RegisterMetaServiceHandlerClient(ctx, mux, metapb.NewMetaServiceClient(metaConn)); err != nil {
		return fmt.Errorf("failed to create a meta grpc client: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to serve http server: %w", err)
	case <-ctx.Done():
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}

		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to close http server: %w", err)
		}

		return nil
	}
}
