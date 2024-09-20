package auth

import (
	"auth/internal/domain/models"

	"auth/internal/services/auth/mocks"
	"auth/internal/storage"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"testing"
)

func Test_Auth_RegisterNewUser(t *testing.T) {
	ctx := context.Background()

	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	tests := []struct {
		nameTest       string
		name           string
		username       string
		password       []byte
		mockUserSaver  func(name, username string, passHash []byte) UserSaver
		expectedID     int
		expectedErrStr string
	}{
		{
			nameTest: "Success",
			name:     "Matvey",
			username: "MatveyTabby",
			password: []byte(mock.AnythingOfType("[]byte")),
			mockUserSaver: func(name, username string, passHash []byte) UserSaver {
				s := mocks.NewUserSaver(t)

				s.EXPECT().
					SaveUser(ctx, name, username, mock.Anything).
					Return(1, nil)
				return s
			},
			expectedID:     1,
			expectedErrStr: "",
		},
		{
			nameTest: "User already exists",
			name:     "Matvey",
			username: "MatveyTabby",
			password: []byte(mock.AnythingOfType("[]byte")),
			mockUserSaver: func(name, username string, passHash []byte) UserSaver {
				s := mocks.NewUserSaver(t)
				s.EXPECT().
					SaveUser(ctx, name, username, mock.Anything).
					Return(0, ErrUserExists)
				return s
			},
			expectedID:     0,
			expectedErrStr: "user already exists",
		},
		{
			nameTest: "Another error during register",
			name:     "Matvey",
			username: "MatveyTabby",
			password: []byte(mock.AnythingOfType("[]byte")),
			mockUserSaver: func(name, username string, passHash []byte) UserSaver {
				s := mocks.NewUserSaver(t)
				s.EXPECT().
					SaveUser(ctx, name, username, mock.Anything).
					Return(0, fmt.Errorf("failed to save user"))
				return s
			},
			expectedID:     0,
			expectedErrStr: "failed to save user",
		},
	}

	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {

			s := Auth{
				UserSaver: tc.mockUserSaver(tc.name, tc.username, tc.password),
				log:       log,
			}

			ID, err := s.RegisterNewUser(ctx, tc.name, tc.username, tc.username)

			if tc.expectedErrStr != "" {
				assert.ErrorContains(t, err, tc.expectedErrStr)
				assert.Equal(t, 0, ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, ID)
			}
		})
	}
}

func Test_Auth_Login(t *testing.T) {
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
		passHash       func(password string) []byte
		mockProvider   func(name, username, passHash string) UserProvider
		Token          string
		expectedErrStr string
	}{
		{
			nameTest: "Success",
			name:     "Matvey",
			username: "MatveyTabby",
			password: "123456",
			mockProvider: func(name, username, password string) UserProvider {
				s := mocks.NewUserProvider(t)

				passHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

				s.EXPECT().User(ctx, username).
					Return(models.User{
						Name:     name,
						Username: username,
						PassHash: passHash,
					}, nil)
				return s
			},
			expectedErrStr: "",
			Token:          "",
		},
		{
			nameTest: "Invalid credentials(incorrect username)",
			username: "MatveyTabby",
			mockProvider: func(name, username, password string) UserProvider {
				s := mocks.NewUserProvider(t)

				s.EXPECT().
					User(ctx, username).
					Return(models.User{}, storage.ErrUserNotFound)

				return s
			},
			expectedErrStr: "invalid credentials",
			Token:          "",
		}, {
			nameTest: "Invalid credentials(incorrect password)",
			name:     "Matvey",
			username: "MatveyTabby",
			password: "123456",
			mockProvider: func(name, username, password string) UserProvider {
				s := mocks.NewUserProvider(t)

				passHash := []byte("incorrect_password")

				s.EXPECT().
					User(ctx, username).
					Return(models.User{
						Name:     name,
						Username: username,
						PassHash: passHash,
					}, nil)
				return s
			},
			expectedErrStr: ErrInvalidCredentials.Error(),
			Token:          "",
		},
		{
			nameTest: "Another error during login",
			username: "MatveyTabby",
			mockProvider: func(name, username, passHash string) UserProvider {
				s := mocks.NewUserProvider(t)

				s.EXPECT().
					User(ctx, username).
					Return(models.User{}, fmt.Errorf("another error"))

				return s
			},
			expectedErrStr: "another error",
			Token:          "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.nameTest, func(t *testing.T) {
			s := Auth{
				UserProvider: tc.mockProvider(tc.name, tc.username, tc.password),
				log:          log,
			}

			token, err := s.Login(ctx, tc.username, tc.password)

			if tc.expectedErrStr != "" {
				assert.ErrorContains(t, err, tc.expectedErrStr)
				assert.Equal(t, "", token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}

}
