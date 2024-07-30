package app

import (
	grpcapp "auth/internal/app/grpc"
	"context"
	"log/slog"
	"time"
)

type App struct { //TODO: некая абстракция над структурой которую мы написали в grpcapp
	GRPCServer *grpcapp.App
}

func NewApp(ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	//TODO: инициализировать хранилище (storage)
	//TODO: инициализировать сервисный слой (auth service)
	//TODO: инициализировать grpc приложение

	grpcServer := grpcapp.NewApp(log, grpcPort)

	return &App{
		GRPCServer: grpcServer,
	}
}
