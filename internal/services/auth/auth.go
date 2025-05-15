package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/finlleyl/gRPC/internal/domain/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
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

type Auth struct {
	log         *zap.SugaredLogger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

func New(log *zap.SugaredLogger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	a.log.Info("Registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Errorw("failed to hash password", "error", err, "password", password)
		return 0, fmt.Errorf("%s: %v", ErrInvalidCredentials, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		a.log.Errorw("failed to save user", "error", err, "email", email)
		return 0, fmt.Errorf("%s: %v", ErrInvalidCredentials, err)
	}

	return id, nil
}
