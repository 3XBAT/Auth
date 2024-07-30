package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"context"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//TODO: инициализировать объект конфига
	//TODO: инициализировать логгер
	//TODO: инициализировать приложение
	//TODO: запустить grpc-сервис

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	application := app.NewApp(context.Background(), log, cfg.GRPC.Port, cfg.TokenTTL)

	application.GRPCServer.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
