package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sso/interanal/config"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func main() {
	ctx := context.Background()
	config := config.MustLoad()
	logger := setupLogger(config.Env)
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"starting application",
		slog.String("with config", fmt.Sprintf("%+v", config)),
	)
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
