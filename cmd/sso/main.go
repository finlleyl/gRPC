package main

import (
	"github.com/finlleyl/gRPC/internal/config"
	"github.com/finlleyl/gRPC/internal/logger"
	"log"
)

func main() {
	cfg := config.MustLoad()
	log.Print()

	zLog, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer zLog.Sync()

	zLog.Infow("Logger initialized")
}
