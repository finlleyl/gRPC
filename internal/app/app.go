package app

import (
	grpcapp "github.com/finlleyl/gRPC/internal/app/grpc"
	"github.com/finlleyl/gRPC/internal/services/auth"
	"github.com/finlleyl/gRPC/internal/storage/sqlite"
	"go.uber.org/zap"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(zap *zap.SugaredLogger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(zap, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(zap, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
