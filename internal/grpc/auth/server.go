package auth

import (
	"auth/internal/services/auth"
	"context"
	"errors"
	authv1 "github.com/3XBAT/protos/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// serverAPI is a structure that handles all incoming requests
type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=Auth --with-expecter=true
type Auth interface {
	Login(ctx context.Context,
		username string,
		password string,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		name string,
		username string,
		password string,
	) (userID int, err error)
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context,
	in *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {

	if err := validateLogin(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, in.GetUsername(), in.GetPassword())

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) RegisterNewUser(ctx context.Context,
	in *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	if err := validateRegister(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.GetName(), in.GetUsername(), in.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authv1.RegisterResponse{
		UserId: int64(userID),
	}, nil
}

func validateLogin(in *authv1.LoginRequest) error {
	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if in.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is empty")
	}

	return nil
}

func validateRegister(in *authv1.RegisterRequest) error {

	if in.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is empty")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if in.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is empty")
	}

	return nil
}
