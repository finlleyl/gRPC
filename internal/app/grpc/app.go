package grpcapp

import (
	"fmt"
	"github.com/finlleyl/gRPC/internal/grpc/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *zap.SugaredLogger
	gRPCServer *grpc.Server
	port       int
}

func New(zap *zap.SugaredLogger, authService auth.Auth, port int) *App {
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()))
	auth.Register(gRPCServer, authService)

	return &App{
		log:        zap,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	a.log.Info("starting gRPC server on port", zap.Int("port", a.port))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
