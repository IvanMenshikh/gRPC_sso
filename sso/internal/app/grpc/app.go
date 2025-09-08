package grpc

import (
	authgrpc "gRPC_sso/sso/internal/grpc/auth"
	"log/slog"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// Конструктор
func New(log *slog.Logger, port int) *App {
	gPRCServer := grpc.NewServer() // подключаем grpc
	authgrpc.Register(gPRCServer)  // подключаем обработчик
	return &App{
		log:        log,
		gRPCServer: gPRCServer,
		port:       port,
	}
}
