package main

import (
	"os"

	"github.com/JanCalebManzano/tag-microservices/pkg/runner"

	"github.com/joho/godotenv"

	"github.com/JanCalebManzano/tag-microservices/services/gateway/http"
)

func main() {
	_ = godotenv.Load(".env")

	os.Exit(
		runner.Run(
			"APPS_GATEWAY_PORT", "gateway", runner.HttpServer, http.RunServer),
	)
}
