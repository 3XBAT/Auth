package auth

//это наши хэндлеры
import (
	"context"
	authv1 "github.com/3XBAT/protos"
	"google.golang.org/grpc"
)

// serverAPI is a structure that handles all incoming requests
type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context,
		username string,
		password string,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		name string,
		username string,
		password string,
	) (userID int64, err error)
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	//TODO
	panic("implement me")
}

func (s *serverAPI) RegisterNewUser(ctx context.Context,
	in authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	//TODO
	panic("implement me")
}
