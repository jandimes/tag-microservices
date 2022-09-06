package interceptor

import (
	"context"

	grpccontext "github.com/JanCalebManzano/tag-microservices/pkg/grpc/context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func NewRequestLogger(logger zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = grpccontext.SetRequestID(ctx)
		reqid := grpccontext.GetRequestID(ctx)

		logger.Info("grpc request",
			zap.String("method", info.FullMethod),
			zap.String("request_id", reqid),
		)

		res, err := handler(ctx, req)

		logger.Info("finished",
			zap.String("method", info.FullMethod),
			zap.Uint32("code", uint32(status.Code(err))),
			zap.String("request_id", reqid),
		)

		return res, err
	}
}
