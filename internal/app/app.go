package app

import (
	grpcapp "auth/internal/app/grpc"
	"auth/internal/config"
	"auth/internal/services/auth"
	"auth/internal/storage"
	"context"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	cfg config.Config,
	tokenTTL time.Duration,
) *App {

	newStorage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	authService := auth.NewAuth(log, newStorage, newStorage, tokenTTL)

	grpcServer := grpcapp.NewApp(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcServer,
	}
}
