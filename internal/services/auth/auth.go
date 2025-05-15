package auth

import (
	"context"
	"errors"
	"github.com/finlleyl/gRPC/internal/domain/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}
