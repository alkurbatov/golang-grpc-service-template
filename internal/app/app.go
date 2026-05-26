// Package app implements application running routine.
package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"google.golang.org/grpc"

	"github.com/alkurbatov/golang-grpc-service-template/internal/config"
	"github.com/alkurbatov/golang-grpc-service-template/internal/controller/grpc/middleware"
	grpcv1 "github.com/alkurbatov/golang-grpc-service-template/internal/controller/grpc/v1"
	httpv1 "github.com/alkurbatov/golang-grpc-service-template/internal/controller/http/v1"
	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/grpcserver"
	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/httpserver"
	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/logging"
)

func Run(cfg *config.Config) error {
	if err := logging.Setup(cfg.LogLevel, cfg.LogJSON); err != nil {
		slog.Default().Error("Failed to configure logger", logging.Err(err))
		return err
	}

	l := slog.Default()
	l.Info(cfg.String())

	grpcSrv := grpcserver.New(
		cfg.GRPCReflection,
		grpc.UnaryInterceptor(middleware.LoggingUnaryInterceptor(l)),
		grpc.StreamInterceptor(middleware.LoggingStreamInterceptor(l)),
	)
	grpcv1.RegisterRoutes(l, grpcSrv.Server)

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	httpSrv := httpserver.NewProdServer(httpv1.RegisterPublicRoutes(reg), cfg.HTTPAddress)

	var profSrv *httpserver.Server
	if cfg.ProfAddress != "" {
		profSrv = httpserver.New(httpv1.RegisterPrivateRoutes(), cfg.ProfAddress)
	}

	err := start(l, cfg, grpcSrv, httpSrv, profSrv)

	l.Info("Stopping the service...")
	shutdown(l, cfg, grpcSrv, httpSrv, profSrv)

	l.Info("Service has stopped")

	return err
}

func start(
	l *slog.Logger,
	cfg *config.Config,
	grpcSrv *grpcserver.Server,
	httpSrv *httpserver.Server,
	profSrv *httpserver.Server,
) error {
	grpcSrvCtx, err := grpcSrv.Start(cfg.GRPCAddress)
	if err != nil {
		l.Error("Failed to start gRPC server", logging.Err(err))
		return err
	}

	httpSrvCtx := httpSrv.Start()

	profSrvCtx := context.Background()
	if profSrv != nil {
		profSrvCtx = profSrv.Start()
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	l.Info("Service has started")

	select {
	case s := <-interrupt:
		l.Info("Received signal " + s.String())

	case <-grpcSrvCtx.Done():
		if err = context.Cause(grpcSrvCtx); err != nil {
			l.Error("gRPC server has failed", logging.Err(err))
		}

	case <-httpSrvCtx.Done():
		err = context.Cause(httpSrvCtx)
		l.Error("HTTP server has failed", logging.Err(err))

	case <-profSrvCtx.Done():
		err = context.Cause(profSrvCtx)
		l.Error("Profiler has failed", logging.Err(err))
	}

	signal.Stop(interrupt)

	return err
}

// shutdown gracefully stops the application.
func shutdown(
	logger *slog.Logger,
	cfg *config.Config,
	grpcSrv *grpcserver.Server,
	httpSrv *httpserver.Server,
	profSrv *httpserver.Server,
) {
	logger.Info("Awaiting prestop timeout...")
	time.Sleep(cfg.PrestopTimeout)

	logger.Info("Shutting down gRPC API...")

	if !grpcSrv.Shutdown(cfg.GRPCTimeout) {
		logger.Warn("Graceful shutdown timeout exceeded, gRPC server was stopped forcibly")
	}

	logger.Info("Shutting down HTTP API...")

	if err := httpSrv.Stop(); err != nil {
		logger.Error("HTTP server shutdown has failed", logging.Err(err))
	}

	if profSrv != nil {
		logger.Info("Shutting down profiler...")

		if err := profSrv.Stop(); err != nil {
			logger.Error("Profiler shutdown has failed", logging.Err(err))
		}
	}
}
