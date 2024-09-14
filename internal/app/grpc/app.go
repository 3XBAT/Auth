package grpcapp

import (
	authgRPC "auth/internal/grpc/auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

// этот файл нужен для того, чтобы разгрузить main
// это приложение, в которое мы оборачиваем наш grpc сервис, для того чтобы было удобнее сконфигурировать его внутри него

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger,
	port int,
	authService authgRPC.Auth) *App {
	grpcServer := grpc.NewServer()

	authgRPC.Register(grpcServer, authService) // подключение обработчика (вроде бы как подрубаем хэндлеры)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}

}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error { // запуск сервера
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op),
		slog.Int("port", a.port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port)) // создаём слушателя порта, который будет ловить
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop() // эта функция заканчивает прием новых запросов, также ждет когда старые обработаются и только потом отрубает приложение
}
