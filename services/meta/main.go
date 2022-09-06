package main

import (
	"os"

	"github.com/JanCalebManzano/tag-microservices/pkg/runner"

	"github.com/joho/godotenv"

	"github.com/JanCalebManzano/tag-microservices/services/meta/grpc"
)

func main() {
	_ = godotenv.Load(".env")

	os.Exit(
		runner.Run(
			"APPS_META_SERVICE_PORT", "meta-service", runner.GrpcServer, grpc.RunServer),
	)
}
