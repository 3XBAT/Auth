package auth

import (
	"auth/internal/domain/models"
	"auth/internal/jwt"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

//go:generate  go run github.com/vektra/mockery/v2@latest --name=UserProvider --with-expecter=true
type UserProvider interface {
	User(ctx context.Context,
		username string,
	) (models.User, error)
}

//go:generate  go run github.com/vektra/mockery/v2@latest --name=UserSaver --with-expecter=true
type UserSaver interface {
	SaveUser(ctx context.Context,
		name string,
		username string,
		PassHash []byte,
	) (uid int, err error)
}

type Auth struct { // Repository
	UserProvider
	UserSaver
	log      *slog.Logger
	TokenTTL time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// NewAuth returns a new instance of the Auth service
func NewAuth(
	log *slog.Logger,
	userProvider UserProvider,
	userSaver UserSaver,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		UserProvider: userProvider,
		UserSaver:    userSaver,
		log:          log,
		TokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {

	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)
	log.Info("attempting to login user")

	user, err := a.UserProvider.User(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", "", err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", "", err.Error())
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, a.TokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", "", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	name string,
	username string,
	pass string,
) (int, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)
	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.UserSaver.SaveUser(ctx, name, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user already exists", "", err.Error())

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists) // сделано специально, чтобы в хэндлеры не пробрасывалась ошибка соля работы с данными
		}

		log.Error("failed to save user", "", err.Error())

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered successfully")

	return id, nil
}
