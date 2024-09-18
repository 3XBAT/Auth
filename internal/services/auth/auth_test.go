package auth

import (
	"auth/internal/domain/models"
	"auth/internal/services/auth/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestAuth_RegisterNewUser(t *testing.T) {
	ctx := context.Background()

	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	tests := []struct {
		nameTest       string
		name           string
		username       string
		password       string
		mockSaver      func(name, password, username string) UserSaver
		expectedID     int
		expectedErrStr string
	}{
		{
			nameTest: "Success",
			name:     "John",
			username: "JohnTravolta",
			password: "123456",
			mockSaver: func(name, username, password string) UserSaver {
				s := mocks.NewUserSaver(t)

				s.EXPECT().
					SaveUser(ctx, name, username, []byte("123456")).
					Return(1, nil)

				return s
			},
			expectedID:     1,
			expectedErrStr: "",
		},
		{
			nameTest: "User already exists",
			name:     "John",
			username: "JohnTravolta",
			password: "123456",
			mockSaver: func(name, username, password string) UserSaver {
				s := mocks.NewUserSaver(t)
				s.EXPECT().SaveUser(ctx, name, username, []byte("123456")).
					Return(0, fmt.Errorf("user already exists"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "user already exists",
		},
		{
			nameTest: "another error during auth",
			name:     "John",
			username: "JohnTravolta",
			password: "123456",
			mockSaver: func(name, username, password string) UserSaver {
				s := mocks.NewUserSaver(t)
				s.EXPECT().SaveUser(ctx, name, username, []byte("123456")).
					Return(0, errors.New("failed to save user"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "failed to save user",
		},
	}

	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {

			s := Auth{
				UserSaver: test.mockSaver(test.name, test.username, test.password),
				TokenTTL:  time.Hour,
				log:       log,
			}

			resp, err := s.UserSaver.SaveUser(ctx, test.name, test.username, []byte(test.password))
			if test.expectedErrStr != "" {
				assert.ErrorContains(t, err, test.expectedErrStr)
				assert.Equal(t, 0, resp)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp)
			}

		})
	}
}

func TestAuth_Login(t *testing.T) {
	ctx := context.Background()

	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	expectedUser := models.User{
		ID:       1,
		Name:     "John",
		Username: "JohnTravolta",
		PassHash: []byte("123456"),
	}

	tests := []struct {
		nameTest       string
		username       string
		mockProvider   func(username string) UserProvider
		ExpectedUser   models.User
		ExpectedErrStr string
	}{
		{
			nameTest: "Success",
			username: "JohnTravolta",
			mockProvider: func(username string) UserProvider {
				s := mocks.NewUserProvider(t)
				s.EXPECT().User(ctx, username).
					Return(expectedUser, nil)
				return s
			},
			ExpectedUser:   expectedUser,
			ExpectedErrStr: "",
		},
		{
			nameTest: "User not found",
			username: "JohnTravolta",
			mockProvider: func(username string) UserProvider {
				s := mocks.NewUserProvider(t)
				s.EXPECT().User(ctx, username).
					Return(models.User{}, fmt.Errorf("user not found"))
				return s
			},
			ExpectedUser:   models.User{},
			ExpectedErrStr: "user not found",
		},
		{
			nameTest: "Invalid credentials",
			username: "JohnTravolta",
			mockProvider: func(username string) UserProvider {
				s := mocks.NewUserProvider(t)
				s.EXPECT().User(ctx, username).
					Return(models.User{}, fmt.Errorf("invalid credentials"))
				return s
			},
			ExpectedUser:   models.User{},
			ExpectedErrStr: "invalid credentials",
		},
		{
			nameTest: "Another error during auth",
			username: "JohnTravolta",
			mockProvider: func(username string) UserProvider {
				s := mocks.NewUserProvider(t)
				s.EXPECT().User(ctx, username).
					Return(models.User{}, fmt.Errorf("failed to get user"))
				return s
			},
			ExpectedUser:   models.User{},
			ExpectedErrStr: "failed to get user",
		},
	}

	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			s := Auth{
				UserProvider: test.mockProvider(test.username),
				TokenTTL:     time.Hour,
				log:          log,
			}

			resp, err := s.UserProvider.User(ctx, test.username)
			if test.ExpectedErrStr != "" {
				assert.ErrorContains(t, err, test.ExpectedErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedUser, resp)
			}
		})
	}
}
