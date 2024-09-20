package auth

import (
	"auth/internal/services/auth"
	"context"
	"fmt"
	"testing"

	"auth/internal/grpc/auth/mocks"
	authv1 "github.com/3XBAT/protos/gen/go"

	"github.com/stretchr/testify/assert"
)

func Test_serverAPI_Register(t *testing.T) {
	//ErrMock := errors.New("mock error")
	ctx := context.Background()

	tests := []struct {
		nameTest       string
		in             *authv1.RegisterRequest
		mockService    func(string, string, string) Auth
		expectedResp   *authv1.RegisterResponse
		expectedErr    error
		expectedErrStr string
	}{
		{
			nameTest: "Success",
			in: &authv1.RegisterRequest{
				Name:     "Matvey",
				Username: "MatveyTabby",
				Password: "OOP",
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().
					RegisterNewUser(ctx, name, username, password).
					Return(1, nil).Once()
				return s
			},
			expectedResp: &authv1.RegisterResponse{
				UserId: 1,
			},
			expectedErr: nil,
		},
		{
			nameTest: "Empty Name",
			in: &authv1.RegisterRequest{
				Name:     "",
				Username: "JohnTravolta",
				Password: "a1b2c5",
			},
			mockService: func(name, username, password string) Auth {
				return nil
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("name is empty"),
		},
		{
			nameTest: "Empty username",
			in: &authv1.RegisterRequest{
				Name:     "Matvey",
				Username: "",
				Password: "OOP",
			},
			mockService: func(name, username, password string) Auth {
				return nil
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("username is empty"),
		},
		{
			nameTest: "Empty Password",
			in: &authv1.RegisterRequest{
				Name:     "Matvey",
				Username: "MatveyTabby",
				Password: "",
			},
			mockService: func(name, username, password string) Auth {
				return nil
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("password is empty"),
		},
		{
			nameTest: "Full Empty",
			in: &authv1.RegisterRequest{
				Name:     "",
				Username: "",
				Password: "",
			},
			mockService: func(name, username, password string) Auth {
				return nil
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("name is empty"),
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
			expectedErr:  fmt.Errorf("password is empty"),
		},
		{
			nameTest: "User already exists",
			in: &authv1.RegisterRequest{
				Name:     "Matvey Tabby",
				Username: "JohnTravolta",
				Password: "a1b2c5",
			},
			mockService: func(name, username, password string) Auth {
				s := mocks.NewAuth(t)

				s.EXPECT().RegisterNewUser(ctx, name, username, password).
					Return(0, auth.ErrUserExists)
				return s
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("user already exists"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			s := &serverAPI{
				auth: tc.mockService(tc.in.GetName(), tc.in.GetUsername(), tc.in.GetPassword()),
			}

			resp, err := s.RegisterNewUser(ctx, tc.in)

			if err != nil {
				assert.Equal(t, tc.expectedResp, resp)
				assert.ErrorContains(t, err, fmt.Sprintf("%s", tc.expectedErr))

			} else {
				assert.Equal(t, tc.expectedResp, resp)
			}

		})
	}
}

func Test_serverAPI_Login(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		nameTest       string
		in             *authv1.LoginRequest
		mockService    func(string, string) Auth
		expectedResp   *authv1.LoginResponse
		expectedErrStr string
	}{
		{
			nameTest: "Success",
			in: &authv1.LoginRequest{
				Username: "MatveyTabby",
				Password: "OOP",
			},
			mockService: func(username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().
					Login(ctx, username, password).
					Return("generatedToken", nil).Once()

				return s
			},
			expectedResp: &authv1.LoginResponse{
				Token: "generatedToken",
			},
			expectedErrStr: "",
		},
		{
			nameTest: "Empty Username",
			in: &authv1.LoginRequest{
				Username: "",
				Password: "OOP",
			},
			mockService: func(username, password string) Auth {
				return nil
			},
			expectedResp:   nil,
			expectedErrStr: "username is empty",
		},
		{
			nameTest: "Empty Password",
			in: &authv1.LoginRequest{
				Username: "MatveyTabby",
				Password: "",
			},
			mockService: func(username, password string) Auth {
				return nil
			},
			expectedResp:   nil,
			expectedErrStr: "password is empty",
		},
		{
			nameTest: "Invalid credentials",
			in: &authv1.LoginRequest{
				Username: "MatveyTabby",
				Password: "OOP",
			},
			mockService: func(username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().
					Login(ctx, username, password).
					Return("", auth.ErrInvalidCredentials).Once()

				return s
			},
			expectedResp:   nil,
			expectedErrStr: "user not found",
		},
		{
			nameTest: "another error during Login",
			in: &authv1.LoginRequest{
				Username: "MatveyTabby",
				Password: "OOP",
			},
			mockService: func(username, password string) Auth {
				s := mocks.NewAuth(t)
				s.EXPECT().
					Login(ctx, username, password).
					Return("", fmt.Errorf("internal server error")).Once()
				return s
			},
			expectedResp:   nil,
			expectedErrStr: "internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			s := &serverAPI{
				auth: tc.mockService(tc.in.GetUsername(), tc.in.GetPassword()),
			}

			resp, err := s.Login(ctx, tc.in)

			if tc.expectedErrStr != "" {
				assert.ErrorContains(t, err, tc.expectedErrStr)
			} else {
				assert.Equal(t, tc.expectedResp, resp)
				assert.Nil(t, err)
			}

		})
	}
}
