package app

import (
	"log/slog"
	grpcapp "sso/interanal/app/grpc"
	"time"
)

type App struct {
	GrpcServer *grpcapp.GrpcApp
}

func New(logger *slog.Logger, port int, db string, tokenTTL time.Duration) *App {
	// TODO: init postgres
	// TODO: init auth service
	grpcApp := grpcapp.New(logger, port)
	return &App{
		GrpcServer: grpcApp,
	}
}
