package main

import (
	"github.com/finlleyl/gRPC/internal/app"
	"github.com/finlleyl/gRPC/internal/config"
	"github.com/finlleyl/gRPC/internal/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	application := app.New(zLog, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go func() { application.GRPCServer.MustRun() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	zLog.Infow("Graceful shutdown complete")
}
