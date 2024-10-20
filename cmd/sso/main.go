package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso/interanal/app"
	"sso/interanal/config"
	"syscall"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	dbCfg := config.MustLoadDBConfig("./config/db.yaml")
	logger := setupLogger(cfg.Env)
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"starting application",
		slog.String("with config", fmt.Sprintf("%+v", cfg)),
	)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)
	grpcApp := app.New(ctx, logger, cfg.GRPC.Port, dbURL, cfg.TokenTTL)
	go grpcApp.GrpcServer.MustRun()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	grpcApp.GrpcServer.Stop()
	grpcApp.Conn.Stop(ctx)
	logger.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return log
}
