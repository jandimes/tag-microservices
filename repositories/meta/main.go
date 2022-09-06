package main

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/JanCalebManzano/tag-microservices/pkg/runner"
	"github.com/JanCalebManzano/tag-microservices/repositories/meta/grpc"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_ = godotenv.Load(".env")

	os.Exit(
		runner.Run(
			"APPS_META_REPOSITORY_PORT", "meta-repository", runner.GrpcServer, grpc.RunServer),
	)
}
