package auth

import (
	"context"
	"errors"
	ssopb "github.com/finlleyl/gRPC/gen/go/sso"
	"github.com/finlleyl/gRPC/internal/services/auth"
	"github.com/finlleyl/gRPC/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssopb.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssopb.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, in *ssopb.LoginRequest) (*ssopb.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing app ID")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssopb.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *ssopb.RegisterRequest) (*ssopb.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "missing email")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "missing password")
	}

	uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &ssopb.RegisterResponse{UserId: uid}, nil
}
