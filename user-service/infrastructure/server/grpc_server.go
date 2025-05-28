package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/tuandq2112/go-microservices/shared/interceptors"
	"github.com/tuandq2112/go-microservices/shared/locale"
	"github.com/tuandq2112/go-microservices/shared/logger"
	pb "github.com/tuandq2112/go-microservices/shared/proto/types/user"
	"github.com/tuandq2112/go-microservices/shared/server"
	"github.com/tuandq2112/go-microservices/user-service/appconfig"
	"github.com/tuandq2112/go-microservices/user-service/internal/interfaces/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Handler *handler.UserHandler
	Logger  logger.Logger
	Locale  *locale.Locale
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.Logger, s.Handler, s.Locale)

	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.Logger)
	g.Add(signalStart, signalStop)

	if err := g.Run(); err != nil {
		s.Logger.Error("Server exited", zap.Error(err))
	}
}

func buildGRPCServer(logger logger.Logger, handler *handler.UserHandler, locale *locale.Locale) (func() error, func(error)) {
	port := "50052"
	if appconfig.Port != "" {
		port = appconfig.Port
	}

	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLocaleInterceptor(locale),
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
		interceptors.UnaryAuthInterceptor(logger, appconfig.JWTSecret),
	}
	streamInts := []grpc.StreamServerInterceptor{
		interceptors.StreamLocaleInterceptor(locale),
		interceptors.StreamLoggerInterceptor(logger),
		interceptors.StreamRecoveryInterceptor(logger),
		interceptors.StreamAuthInterceptor(logger, appconfig.JWTSecret),
	}

	addr := fmt.Sprintf(":%s", port)
	logger.Info("Starting gRPC server", zap.String("address", addr))

	_, start, stop, err := server.BuildGRPCServer(
		addr,
		func(s *grpc.Server) {
			logger.Info("Registering UserServiceServer")
			pb.RegisterUserServiceServer(s, handler)
		},
		unaryInts,
		streamInts,
	)
	if err != nil {
		logger.Error("Failed to build gRPC server", zap.Error(err))
		return nil, nil
	}

	return start, stop
}

func buildSignalHandler(logger logger.Logger) (func() error, func(error)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return func() error {
			sig := <-sigs
			return fmt.Errorf("signal received: %v", sig)
		}, func(err error) {
			close(sigs)
		}
}
