package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/JanCalebManzano/tag-microservices/pkg/logger"
	"go.uber.org/zap"
)

type runService func(ctx context.Context, port int, logger *zap.Logger) error

type serverType string

const (
	HttpServer serverType = "http"
	GrpcServer serverType = "grpc"
)

func Run(portEnv, loggerName string, svrType serverType, run runService) int {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	l, err := logger.New()
	if err != nil {
		_, ferr := fmt.Fprintf(os.Stderr, "failed to create logger: %s", err)
		if ferr != nil {
			panic(fmt.Sprintf("failed to write log:`%s` original error is:`%s`", ferr, err))
		}
		return 1
	}
	lNamed := l.Named(loggerName)

	portRaw := os.Getenv(portEnv)
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		port = 0
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- run(ctx, port, lNamed.Named(string(svrType)))
	}()

	lNamed.Info("Started server...", zap.Int("port", port))

	select {
	case err := <-errCh:
		lNamed.Error(err.Error())
		return 1
	case <-ctx.Done():
		lNamed.Info("shutting down...")
		return 0
	}
}
