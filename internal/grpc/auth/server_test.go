package auth

import (
	"context"
	"errors"
	"testing"

	"auth/internal/grpc/auth/mocks"
	authv1 "github.com/3XBAT/protos/gen/go"

	"github.com/stretchr/testify/assert"
)

func Test_serverAPI_Register(t *testing.T) {
	ErrMock := errors.New("mock error")
	ctx := context.Background()

	tests := []struct {
		nameTest                    string
		in                          *authv1.RegisterRequest
		name                        string
		username                    string
		password                    string
		mockUnimplementedAuthServer func() authv1.UnimplementedAuthServer
		mockService                 func(string, string, string) Auth
		expectedID                  int
		expectedResp                *authv1.RegisterResponse
		expectedErr                 error
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
			in: &authv1.RegisterRequest{
				Name:     "",
				Username: "JohnTravolta",
				Password: "a1b2c5",
			},
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
		{
			nameTest: "error during RegisterNewUser",
			in: &authv1.RegisterRequest{
				Name:     "matvey Tabby",
				Username: "JohnTravolta",
				Password: "",
			},
			mockService: func(name, username, password string) Auth {
				return nil
			},
			expectedResp: nil,
			expectedErr:  ErrMock,
		},
	}
	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			s := &serverAPI{
				auth: tc.mockService(tc.in.GetName(), tc.in.GetUsername(), tc.in.GetPassword()),
			}

			resp, _ := s.RegisterNewUser(ctx, tc.in)

			assert.Equal(t, tc.expectedResp, resp)
			//assert.ErrorIs(t, err, tc.expectedErr)

			//if tc.expectedErrStr == "" {
			//	assert.Nil(t, err)
			//} else {
			//	assert.Error(t, err)
			//	assert.Equal(t, tc.expectedID, resp)
			//}
			//
			//assert.Equal(t, tc.expectedID, resp)
		})
	}
}
