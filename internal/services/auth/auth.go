package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/finlleyl/gRPC/internal/domain/models"
	"github.com/finlleyl/gRPC/internal/lib/jwt"
	"github.com/finlleyl/gRPC/internal/storage"
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

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	a.log.Info("Logging in")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")

			return "", ErrInvalidCredentials
		}

		a.log.Errorw("failed to get user", "error", err, "email", email)
		return "", fmt.Errorf("%s: %v", ErrInvalidCredentials, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %v", ErrInvalidCredentials, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %v", ErrInvalidCredentials, err)
	}

	a.log.Info("Logged in", "email", email)

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Errorw("failed to create token", "error", err, "email", email)
		return "", fmt.Errorf("%s: %v", ErrInvalidCredentials, err)
	}

	return token, nil
}
