package app

import (
	"context"
	"log/slog"
	grpcapp "sso/interanal/app/grpc"
	"sso/interanal/service/auth"
	"sso/interanal/storage/postgres"
	"time"
)

type App struct {
	GrpcServer *grpcapp.GrpcApp
	Conn       *postgres.Storage
}

func New(ctx context.Context, logger *slog.Logger, port int, db string, tokenTTL time.Duration) *App {
	storage := postgres.MustNewConnection(ctx, db)
	authService := auth.New(logger, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(logger, authService, port)
	return &App{
		GrpcServer: grpcApp,
		Conn:       storage,
	}
}
