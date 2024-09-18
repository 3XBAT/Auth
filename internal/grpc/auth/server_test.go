package auth

import (
	"auth/internal/grpc/auth/mocks"
	"context"
	"errors"
	authv1 "github.com/3XBAT/protos/gen/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_serverAPI_Register(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		nameTest                    string
		name                        string
		username                    string
		password                    string
		mockUnimplementedAuthServer func() authv1.UnimplementedAuthServer
		mockService                 func(string, string, string) Auth
		expectedID                  int
		expectedErrStr              string
	}{
		{
			nameTest: "Success",
			name:     "John",
			username: "JohnTravolta",
			password: "a1b2c5",
			mockUnimplementedAuthServer: func() authv1.UnimplementedAuthServer {
				return authv1.UnimplementedAuthServer{}
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().
					RegisterNewUser(ctx, name, username, password).
					Return(1, nil)
				return s
			},
			expectedID:     1,
			expectedErrStr: "",
		},
		{
			nameTest: "Empty Name",
			name:     "",
			username: "JohnTravolta",
			password: "a1b2c5",
			mockUnimplementedAuthServer: func() authv1.UnimplementedAuthServer {
				return authv1.UnimplementedAuthServer{}
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().RegisterNewUser(ctx, name, username, password).
					Return(0, errors.New("name is empty"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "name is empty",
		},
		{
			nameTest: "Empty username",
			name:     "John",
			username: "",
			password: "a1b2c5",
			mockUnimplementedAuthServer: func() authv1.UnimplementedAuthServer {
				return authv1.UnimplementedAuthServer{}
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().RegisterNewUser(ctx, name, username, password).
					Return(0, errors.New("username is empty"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "username is empty",
		},
		{
			nameTest: "Empty Password",
			name:     "John",
			username: "JohnTravolta",
			password: "",
			mockUnimplementedAuthServer: func() authv1.UnimplementedAuthServer {
				return authv1.UnimplementedAuthServer{}
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().RegisterNewUser(ctx, name, username, password).
					Return(0, errors.New("password is empty"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "password is empty",
		},
		{
			nameTest: "Full Empty",
			name:     "",
			username: "",
			password: "",
			mockUnimplementedAuthServer: func() authv1.UnimplementedAuthServer {
				return authv1.UnimplementedAuthServer{}
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().RegisterNewUser(ctx, name, username, password).
					Return(0, errors.New("name is empty"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "name is empty",
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			s := &serverAPI{
				UnimplementedAuthServer: test.mockUnimplementedAuthServer(),
				auth:                    test.mockService(test.name, test.username, test.password),
			}

			resp, err := s.auth.RegisterNewUser(ctx, test.name, test.username, test.password)

			if test.expectedErrStr == "" {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedID, resp)
			}

			assert.Equal(t, test.expectedID, resp)
		})
	}
}
