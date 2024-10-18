package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authgrpc "sso/interanal/grpc/auth"

	"google.golang.org/grpc"
)

type GrpcApp struct {
	logger     *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int) *GrpcApp {
	grpcServer := grpc.NewServer()
	authgrpc.RegisterServerAPI(grpcServer)
	return &GrpcApp{
		logger:     logger,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *GrpcApp) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *GrpcApp) run() error {
	const op = "grpcapp.Run"
	logger := a.logger.With(slog.String("op", op))
	logger.Info("starting grpc server")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("gepc server is running", slog.String("port", l.Addr().String()))

	if err := a.grpcServer.Serve(l); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *GrpcApp) Stop() {
	const op = "grpcapp.Stop"
	a.logger.Info("stopping server", slog.String("op", op), slog.Int("port", a.port))
	a.grpcServer.GracefulStop()
}
