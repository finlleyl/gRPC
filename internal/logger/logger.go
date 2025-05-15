package logger

import (
	"fmt"
	"github.com/finlleyl/gRPC/internal/config"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	switch cfg.Env {
	case "local":
		logger, err = zap.NewDevelopment()
	case "dev", "prod":
		logger, err = zap.NewProduction()
	default:
		return nil, fmt.Errorf("unknown logger enviroment: %s", cfg.Env)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger.Sugar(), nil
}
